<script>
	import Router from 'svelte-spa-router'

	import Home from "./Home.svelte";
	import Worker from "./Worker.svelte";

	const routes = {
		'/': Home,
		'/workers/:id': Worker,
	}

</script>

<body>
	<Router {routes}/>
</body>