import svelte from "rollup-plugin-svelte";
import resolve from "rollup-plugin-node-resolve";
import commonjs from "rollup-plugin-commonjs";
import postcss from "rollup-plugin-postcss";
import { terser } from "rollup-plugin-terser";

export default {
  input: "src/main.js",
  output: {
    // sourcemap: true,
    format: "iife",
    name: "app",
    file: "public/assets/bundle.js"
  },
  plugins: [
    postcss({
      extensions: [".css"]
    }),
    svelte({
      dev: false
      //   css: css => css.write("public/assets/bundle.css")
    }),
    resolve({
      browser: true,
      dedupe: importee =>
        importee === "svelte" || importee.startsWith("svelte/")
    }),
    commonjs(),
    terser()
  ]
};
