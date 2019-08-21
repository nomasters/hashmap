import App from "./App.svelte";
import "../node_modules/milligram/dist/milligram.min.css";
import "./custom.css";

const app = new App({
  target: document.body,
  props: {
    name: "nathan"
  }
});

export default app;
