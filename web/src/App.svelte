<script>
	let workers = [];

	async function getWorkers() {
		const res = await fetch("/api/workers");

		if (res.ok) {
			workers = await res.json();
		} else {
			throw new Error(res.status);
		}
	}
	
	setInterval(getWorkers(), 2000);
</script>

<div id="workersDiv">
	<h2>Workers</h2>

	<ul>
		{#each workers as { ContainerId, Package }}
			<p>Package: <b>{Package}</b></p>
			<p>ID: <b>{ContainerId}</b></p>
		{/each}
	</ul>
</div>

<style>
	:global(body) {
		background-color: rgb(30, 30, 30);
	}

	#workersDiv {
		background-color: rgb(175, 175, 175);
		color: black;

		border-radius: 10px;

		padding-left: 20px;
		padding-right: 20px;
		padding-top: 10px;
		padding-bottom: 10px;
	}

	#workersDiv h2 {
		text-align: center;
	}

	#workersDiv ul {
		background-color: rgb(42, 105, 199);
		border-radius: 5px;
		margin: 10px;
		padding: 15px;
	}
</style>
