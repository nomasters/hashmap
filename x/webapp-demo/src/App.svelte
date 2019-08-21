<script>
  import { Router, Link, Route } from "svelte-routing";
  import WasmExec from "./WasmExec.svelte";
  import Logo from "./Logo.svelte";
  import Get from "./routes/Get.svelte";
  import Home from "./routes/Home.svelte";
  export let name;
  export let url = "";
</script>

<style>

</style>

<Logo width="190px" />
<WasmExec />

<h1>Hello {name}!</h1>

<Router {url}>
  <nav>
    <Link to="/">Home</Link>
    <Link to="/fixedHash">test hash</Link>
  </nav>
  <div>
    <Route path="/:hash" component={Get} />
    <Route path="/" component={Home} />
  </div>
</Router>
