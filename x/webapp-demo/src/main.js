import "./wasm_exec.js";
import "../node_modules/normalize.css/normalize.css";
import "../node_modules/milligram/dist/milligram.min.css";
import "./custom.css";
import App from "./App.svelte";

const app = new App({
  target: document.body,
  props: {
    name: "nathan"
  }
});

export default app;
