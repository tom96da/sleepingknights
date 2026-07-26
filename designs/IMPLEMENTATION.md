# SleepingKnights 実装設計書

**対象:** Lexer, Parser, Evaluator (Go実装)
**Phase:** 1 (Go製インタプリタ基盤)
**最終更新:** 2026-07-14

---

このドキュメントは言語仕様を実装する際の技術的な設計を記載します。
言語そのもの定義については [LANGUAGE_SPEC.md](LANGUAGE_SPEC.md) を参照してください。

---

## 1. 字句解析 (Lexer)

### 1.1 トークン種別

| 種別 | 例 | 正規表現 |
|------|------|--------|
| **キーワード** | `fn`, `if`, `else`, `for`, `while`, `return`, `true`, `false` | 予約語 |
| **識別子** | `x`, `myFunc`, `_private` | `[a-zA-Z_][a-zA-Z0-9_]*` |
| **整数リテラル** | `42`, `0x1F`, `0b1010` | `[0-9]+` \| `0x[0-9a-fA-F]+` \| `0b[01]+` |
| **浮動小数点リテラル** | `3.14`, `1.0e-3` | `[0-9]+\.[0-9]+` \| `[0-9]*\.[0-9]+([eE][+-]?[0-9]+)?` |
| **文字列リテラル** | `"hello"`, `"hello\nworld"` | `"([^"\\]|\\.)*"` |
| **演算子** | `+`, `-`, `*`, `/`, `==`, `!=`, `and`, `or`, `not` 等 | 多文字マッチ対応 |
| **区切り文字** | `(`, `)`, `{`, `}`, `,`, `;`, `:`, `=` | 単一文字 |
| **コメント** | `// line`, `/* block */`, `/** doc */` | ツール向けにトークン化し、通常コメントは構文木に含めない |

### 1.2 予約語リスト

```
fn let const if else for while return true false break continue int float bool string void and or not
```

### 1.3 リテラル処理

#### 整数
```
10進:   [0-9]+
16進:   0x[0-9a-fA-F]+
2進:    0b[01]+
```

#### 浮動小数点
```
通常表記: [0-9]+\.[0-9]+
指数表記: [0-9]*\.[0-9]+([eE][+-]?[0-9]+)?
```

#### 文字列
```
エスケープシーケンス: \\, \", \n, \t, \r, \b, \f
```

### 1.4 コメント処理

- `// ...` は行コメント `LINE_COMMENT` としてトークン化する
- `/* ... */` はブロックコメント `BLOCK_COMMENT` としてトークン化する
- `/** ... */` はドキュメントコメント `DOC_COMMENT` としてトークン化する
- ドキュメントコメントは直後の宣言に関連付けるため保持する
- 通常の行コメントとブロックコメントは、処理系本体の parser / evaluator では無視する
  リンターやフォーマッタはコメントトークンを利用できる前提とする

---

## 2. 構文解析 (Parser)

### 2.1 優先度ベースの式パーサ

| 優先度 | 演算子 | 結合性 | 説明 |
|--------|--------|--------|------|
| 1 (最高) | `()` | 左 | 関数呼び出し、グループ化 |
| 2 | `not` | 右 | 論理否定 |
| 3 | `*`, `/`, `%` | 左 | 乗除 |
| 4 | `+`, `-` | 左 | 加減 |
| 5 | `<`, `<=`, `>`, `>=` | 左 | 比較 |
| 6 | `==`, `!=` | 左 | 等価 |
| 7 | `and` | 左 | 論理積（短絡評価） |
| 8 (最低) | `or` | 左 | 論理和（短絡評価） |

### 2.2 文法（BNF形式）

**注:**
- ステートメント終端は **NEWLINE** で判定
- `for` ループの初期化・条件・更新部分では `;` で区切る（括弧内）
- 代入は式の最も低い優先度で、右結合
- 例：`a = b = c` は `a = (b = c)` と評価

```
program        := function_decl* statement*

function_decl  := "fn" identifier "(" parameters? ")" (":" type)? block
parameters     := parameter ("," parameter)*
parameter      := identifier ":" type
type           := "int" | "float" | "bool" | "string" | "void"

block          := "{" statement* "}"

statement      := var_decl
               | const_decl
               | if_stmt
               | while_stmt
               | for_stmt
               | return_stmt
               | break_stmt
               | continue_stmt
               | expression

var_decl       := "let" identifier "=" expression
const_decl     := "const" identifier "=" expression

if_stmt        := "if" "(" if_condition ")" block else_clause?
if_condition   := expression
               | "let" identifier "=" expression
else_clause    := "else" (if_stmt | block)

while_stmt     := "while" "(" expression ")" block

for_stmt       := "for" "(" "let" identifier "=" expression ";" expression ";" expression ")" block

return_stmt    := "return" expression?
break_stmt     := "break"
continue_stmt  := "continue"

expression     := assignment
assignment     := identifier "=" assignment | logical_or
logical_or     := logical_and ( "or" logical_and )*
logical_and    := equality ( "and" equality )*
equality       := comparison ( ("==" | "!=") comparison )*
comparison     := additive ( ("<" | "<=" | ">" | ">=") additive )*
additive       := multiplicative ( ("+" | "-") multiplicative )*
multiplicative := unary ( ("*" | "/" | "%") unary )*
unary          := ("not" | "-")? primary
primary        := INTEGER | STRING | BOOL | identifier |
                  "(" expression ")" | function_call

function_call  := identifier "(" arguments? ")"
arguments      := expression ("," expression)*
```

---

## 3. 抽象構文木 (AST)

### 3.1 ノード型の定義

AST はすべてのノードが `Node` インターフェースを実装する必要があります。各ノードは行・列情報を保持します。

**定義が必要なノード型：**

**プログラム構造:**
- `Program`: 全体のプログラムノード
- `Function`: 関数定義（名前、パラメータ、戻り値型、本体）
- `Parameter`: 関数パラメータ（名前、型）
- `Type`: 型情報（"int", "float", "bool", "string", "void"）

**文:**
- `Block`: ブロック（複数の文）
- `VariableDecl`: 変数宣言（`let x = expr`）
- `Assignment`: 代入（`x = expr`）
- `IfStatement`: if-else 文（条件、then ブロック、else ブロック）
- `WhileStatement`: while ループ（条件、本体）
- `ForStatement`: for ループ（初期化、条件、更新、本体）
- `ReturnStatement`: return 文（戻り値式）
- `BreakStatement`: break 文
- `ContinueStatement`: continue 文
- `ExpressionStatement`: 式文（ステートメント終端は NEWLINE）

**式:**
- `BinaryExpression`: 二項演算（左、演算子、右）
- `UnaryExpression`: 単項演算（演算子、オペランド）
- `FunctionCall`: 関数呼び出し（関数名、引数リスト）
- `Identifier`: 識別子（名前）

**リテラル:**
- `IntegerLiteral`: 整数リテラル（10進、16進、2進）
- `FloatLiteral`: 浮動小数点リテラル
- `StringLiteral`: 文字列リテラル
- `BooleanLiteral`: 真偽値リテラル（true/false）

すべてのノードは **行番号と列番号** を記録して、エラー報告の精度を高めます。

---

## 4. 評価器 (Evaluator)

### 4.1 環境管理

**構造:**
- 階層的な環境（グローバル環境 → 関数スコープ → ブロックスコープ）
- 各環境は以下を管理する必要があります：
  - `parent`: 親環境へのポインタ（グローバルは nil）
  - `variables`: 変数名から値へのマッピング
  - `functions`: 関数名から関数定義へのマッピング

**機能:**
- 変数定義・参照・更新
- 変数のシャドーイング（同じ名前で異なるスコープの変数を区別）
- 関数定義・呼び出し

### 4.2 式評価

各式ノード型を評価して、対応する Go の値を返します：

- **リテラル**: 値をそのまま返す
- **Identifier**: 環境から変数値を参照
- **BinaryExpression**: 左右を評価し、演算子に応じた演算を実行（短絡評価対応）
- **UnaryExpression**: オペランドを評価し、単項演算を実行
- **FunctionCall**: 関数を環境から取得し、引数を評価して呼び出し

### 4.3 制御フロー

**実装方針:**
- **return:** 特殊な制御フロー値を使ってスタック巻き戻しを実現
- **break/continue:** ループ脱出・継続用のフラグを用いた処理
- **スコープ管理:** 関数/ブロック進入時に新環境を作成、離脱時に親環境に戻す

---

## 5. 組み込み関数

Phase 0では基本的な組み込み関数のみ提供します。Phase 3以降で段階的に拡張予定です。

### 5.1 `print` 関数

**機能:**
- 任意の値をコンソールに出力
- 最後に改行を出力

**シグネチャ:** `print(value)`

**評価器での実装方針:**
- 組み込み関数として環境に登録
- 評価時に FunctionCall ノードから呼び出される

---

## 6. エラーハンドリング

### 6.1 エラー種別

| 種別 | 例 |
|------|------|
| **構文エラー** | 予期しないトークン |
| **実行時エラー** | 未定義変数へのアクセス |
| **型エラー** (Phase 3) | 型の不一致 |

### 6.2 エラー出力形式

```
src/program.slk:5:10: error: undefined variable 'x'
  | print(x)
  |       ^
```

- ファイル名、行番号、列番号を記録
- エラー行と問題箇所のポインタを表示

---

## 7. テスト戦略

### 7.1 Lexer テスト

**テスト項目:**
- キーワード認識（fn, let, if, else など）
- 識別子とリテラル分類
- 演算子・区切り文字の正確な分割
- コメント削除
- 複数行入力での行・列情報の正確さ
- エスケープシーケンス処理
- エラーケース（不正なリテラルなど）

### 7.2 Parser テスト

**テスト項目:**
- 式の優先度が正しく反映されているか
- 各種文（if, while, for, return など）の正しい解析
- 関数定義と呼び出しの解析
- エラー位置の正確さ

### 7.3 Evaluator テスト

**テスト項目:**
- リテラル評価
- 算術演算・比較演算・論理演算
- 変数の定義・参照・再代入
- シャドーイング
- 関数定義と呼び出し
- 制御フロー（if, while, for, break, continue, return）
- スコープ管理
- print 関数
- エラーケース（未定義変数など）

---

## 8. 実装チェックリスト

### Lexer
- [ ] トークン分類と正規表現マッチ
- [ ] コメント削除
- [ ] 行・列情報の正確な追跡
- [ ] エスケープシーケンス処理

### Parser
- [ ] 優先度ベースの式パーサ実装
- [ ] 文パーサ実装
- [ ] 関数宣言パーサ実装
- [ ] エラー回復（可能な限り）

### Evaluator
- [ ] 環境（スコープ）管理
- [ ] 式評価ロジック
- [ ] 文実行ロジック
- [ ] 制御フロー（if/for/while/return）
- [ ] 組み込み関数の実装

---

## 9. 今後のフェーズ向け考慮事項

### Phase 3: 型推論
- AST に型情報フィールドを追加可能な設計
- 型チェッカーを Evaluator の前段に挿入

### Phase 4: OOP
- オブジェクト（構造体）を環境に保持
- メソッド呼び出し構文を Parser に追加

### Phase 5: 関数型
- 関数値（クロージャ）を環境に保持
- 高階関数対応

### Phase 6: LLVM
- AST を LLVM IR に変換するコンパイラ層
- 現在の Evaluator とのセマンティクス検証

