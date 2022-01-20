<script>
    import { link } from "svelte-spa-router";

    let workers = [];
    let packages = [];
    let addingPackage = false;
    let serverStatus = "";

    async function addPackage() {
        const pkg = prompt("Enter AUR package name");

        if (pkg != "") {
            addingPackage = true;

            const response = await fetch("/api/packages", {
                method: "POST",
                body: pkg,
            });

            if (!response.ok) {
                alert("Failed to add package: " + response.statusText);
            }

            addingPackage = false;
        }
    }

    async function getWorkers() {
        const res = await fetch("/api/workers");

        if (res.ok) {
            const workersNew = await res.json();

            if (workersNew.toString() != workers.toString()) {
                workers = workersNew;

                getPackages();
            }
        } else {
            throw new Error(res.status);
        }
    }

    async function getPackages() {
        const res = await fetch("/api/packages");

        if (res.ok) {
            packages = await res.json();
        } else {
            throw new Error(res.status);
        }
    }

    async function getStatus() {
        const res = await fetch("/api/status");

        if (res.ok) {
            let statusList = await res.json();

            if (statusList.length == 0) {
                serverStatus = "Idle";
            } else {
                serverStatus = statusList.join(", ");
            }
        } else {
            throw new Error(res.status);
        }
    }

    async function startUpdate() {
        await fetch("/api/update", {
            method: "POST",
        });
    }

    async function deletePackage(pkgName) {
        if (confirm("Are you sure you want to delete " + pkgName + "?")) {
            await fetch("/api/packages", {
                method: "DELETE",
                body: pkgName,
            });

            getPackages();
        }
    }

    getPackages();
    getWorkers();
    getStatus();

    setInterval(getWorkers, 2000);
    setInterval(getStatus, 1000);

    let repoLink =
        window.location.protocol + "//" + window.location.host + "/repo";
</script>

<div>
    <div class="section">
        <h1>Status: <b>{serverStatus}</b></h1>
    </div>
    <div class="section">
        <button on:click={startUpdate}>Check for updates</button>
    </div>
    <div id="top-section">
        <div class="section">
            <h2>Repository</h2>
            <h3>Pacman config snippet:</h3>
            <code>
                [aurer]<br />
                Server = {repoLink}
            </code><br />
            <a href="{repoLink}/">Browse online</a>
        </div>
        <div class="section">
            <h2>Packages</h2>

            <table>
                <tr>
                    <th>Name</th>
                    <th>Version</th>
                    <th>Download</th>
                    <th>Delete</th>
                </tr>
                {#each packages as Package}
                    <tr>
                        <td>{Package.Name}</td>
                        <td>{Package.Version}</td>
                        <td><a href="/repo/{Package.Filename}">Link</a></td>
                        <td
                            ><button
                                on:click={() => deletePackage(Package.Name)}
                                >Delete</button
                            ></td
                        >
                    </tr>
                {/each}
            </table>
            {#if addingPackage}
                <p>Adding package...</p>
            {:else}
                <button id="addPackageButton" on:click={addPackage}>
                    Add package
                </button>
            {/if}
        </div>
    </div>
    <div class="section">
        <h2>Workers</h2>

        <ul class="workers">
            {#each workers as { ContainerId, Package }}
                <a href="/workers/{ContainerId}" use:link>
                    <p>Package: <b>{Package}</b></p>
                    <p>ID: <b>{ContainerId}</b></p>
                </a>
            {/each}
        </ul>
    </div>
</div>

<style>
    :global(body) {
        background-color: rgb(30, 30, 30);
    }

    code {
        font-size: large;
    }

    #top-section {
        display: flex;
        flex-wrap: wrap;
    }

    #top-section .section {
        flex: 1;
    }

    .section {
        background-color: rgb(175, 175, 175);
        color: black;

        border-radius: 10px;

        padding-left: 20px;
        padding-right: 20px;
        padding-top: 5px;
        padding-bottom: 10px;

        margin: 20px;
    }

    .section h2 {
        text-align: center;
    }

    .section ul {
        background-color: rgb(42, 105, 199);
        border-radius: 5px;
        margin: 10px;
        padding: 15px;
    }

    .section table,
    th,
    td {
        border: 1px solid black;
        border-collapse: collapse;
    }
    .section table {
        width: 100%;
    }

    .section th,
    td {
        text-align: center;
        padding: 5px;
    }

    .workers a {
        color: black;
    }
</style>
