import {Monaco} from "@monaco-editor/react";
import {languages} from "monaco-editor";

export const CADENCE_LANGUAGE_ID = "cadence";

export default function configureCadence(monaco: Monaco) {
    monaco.languages.register({
      id: CADENCE_LANGUAGE_ID,
      extensions: [".cdc"],
      aliases: ["CDC", "cdc"],
    });

    const languageDef: languages.IMonarchLanguage = {
      keywords: [
        "if", "else", "return", "continue", "break", "while", "pre", "post",
        "prepare", "execute", "import", "from", "create", "destroy", "priv",
        "pub", "get", "set", "log", "emit", "event", "init", "struct",
        "interface", "fun", "let", "var", "resource", "access", "all",
        "contract", "self", "transaction"
      ],
      typeKeywords: [
        "AnyStruct", "AnyResource", "Void", "Never", "String", "Character",
        "Bool", "Self", "Int8", "Int16", "Int32", "Int64", "Int128", "Int256",
        "UInt8", "UInt16", "UInt32", "UInt64", "UInt128", "UInt256", "Word8",
        "Word16", "Word32", "Word64", "Word128", "Word256", "Fix64", "UFix64"
      ],
      operators: [
        "<-", "<=", ">=", "==", "!=", "+", "-", "*", "/", "%", "&", "!",
        "&&", "||", "?", "??", ":", "=", "@"
      ],
      symbols: /[=><!~?:&|+\-*\/\^%]+/,
      escapes: /\\(?:[abfnrtv\\"]|x[0-9A-Fa-f]{1,4}|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/,
      digits: /\d+(_+\d+)*/,
      octaldigits: /[0-7]+(_+[0-7]+)*/,
      binarydigits: /[0-1]+(_+[0-1]+)*/,
      hexdigits: /[[0-9a-fA-F]+(_+[0-9a-fA-F]+)*/,
      tokenizer: {
        root: [[/[{}]/, "delimiter.bracket"], { include: "common" }],
        common: [
          [/[a-z_$][\w$]*/, {
            cases: {
              "@typeKeywords": "keyword",
              "@keywords": "keyword",
              "@default": "identifier"
            }
          }],
          [/[A-Z][\w]*/, "type.identifier"],
          { include: "@whitespace" },
          [/[()\[\]]/, "@brackets"],
          [/[<>](?!@symbols)/, "@brackets"],
          [/@symbols/, {
            cases: {
              "@operators": "delimiter",
              "@default": ""
            }
          }],
          [/(@digits)[eE]([\-+]?(@digits))?/, "number.float"],
          [/(@digits)\.(@digits)([eE][\-+]?(@digits))?/, "number.float"],
          [/0[xX](@hexdigits)/, "number.hex"],
          [/0[oO]?(@octaldigits)/, "number.octal"],
          [/0[bB](@binarydigits)/, "number.binary"],
          [/(@digits)/, "number"],
          [/[;,.]/, "delimiter"],
          [/"([^"\\]|\\.)*$/, "string.invalid"],
          [/"/, "string", "@string_double"],
        ],
        whitespace: [
          [/[ \t\r\n]+/, ""],
          [/\/\*/, "comment", "@comment"],
          [/\/\/.*$/, "comment"]
        ],
        comment: [
          [/[^\/*]+/, "comment"],
          [/\*\//, "comment", "@pop"],
          [/[\/*]/, "comment"]
        ],
        string_double: [
          [/[^\\"]+/, "string"],
          [/@escapes/, "string.escape"],
          [/\\./, "string.escape.invalid"],
          [/"/, "string", "@pop"]
        ],
      },
    };

    monaco.languages.setMonarchTokensProvider(CADENCE_LANGUAGE_ID, languageDef);
  }
