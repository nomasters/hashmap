<script>
  import "./wasm_exec.js";

  if (!WebAssembly.instantiateStreaming) {
    // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const go = new Go();
  let mod, inst;
  WebAssembly.instantiateStreaming(
    fetch("assets/hashmap.wasm"),
    go.importObject
  )
    .then(result => {
      mod = result.module;
      inst = result.instance;
      document.getElementById("runButton").disabled = false;
    })
    .catch(err => {
      console.error(err);
    });

  async function run() {
    console.clear();
    await go.run(inst);
    inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
  }
</script>

<div>
  <button class="button button-outline" on:click={run} id="runButton" disabled>
    Outlined Button
  </button>
</div>
