// probe-r5.mjs — renderer 0.3.x external 化产物断言（VP-024 R5 · 判据 #5）。
// 检查安装后的 @magicvr/schema-ui-renderer：
//   1) index.js 含 >=3 处 from "@magicvr/schema-ui-..."（依赖图 external 化 · 发布实态全名）
//   2) 无 from "@/..." 残留（内部 alias 清零）
//   3) renderer/index.d.ts 导出面指向包子路径（@magicvr/schema-ui-lib/i18n/runtime）
import { readFileSync } from "node:fs";
import path from "node:path";
import assert from "node:assert";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));
const rendererRoot = path.join(here, "node_modules/@magicvr/schema-ui-renderer");
const js = readFileSync(path.join(rendererRoot, "index.js"), "utf8");

const imports = [...js.matchAll(/from\s+"(@magicvr\/schema-ui-[^"]+)"/g)].map((m) => m[1]);
assert.ok(imports.length >= 3, `external 化断言失败：imports = ${imports.length}`);
assert.ok(imports.every((i) => i.startsWith("@magicvr/schema-ui-")), "全部 import 应为发布实态包子路径");
assert.ok(!js.includes('from "@/'), "index.js 不应含内部 alias 残留");

const dts = readFileSync(path.join(rendererRoot, "renderer/index.d.ts"), "utf8");
assert.ok(dts.includes("@magicvr/schema-ui-lib/i18n/runtime"), "d.ts 导出面应指向包子路径");
assert.ok(!dts.includes('from "@/'), "d.ts 不应含内部 alias 残留");

console.log(`renderer ${JSON.parse(readFileSync(path.join(rendererRoot, "package.json"), "utf8")).version} external 化断言 PASS · imports=${imports.length}（${[...new Set(imports.map((i) => i.split("/")[1]))].join(", ")}）`);