<script>
    import { link } from "svelte-spa-router";

    let workers = [];
    let packages = [];

    async function getWorkers() {
        const res = await fetch("/api/workers");

        if (res.ok) {
            workers = await res.json();
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
    
    getPackages();

    getWorkers();
    setInterval(getWorkers, 2000);

    let repoLink =
        window.location.protocol + "//" + window.location.host + "/repo";
</script>

<div>
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
                </tr>
                {#each packages as Package}
                    <tr>
                        <td>{Package.Name}</td>
                        <td>{Package.Version}</td>
                        <td><a href="/repo/{Package.Filename}">Link</a></td>
                    </tr>
                {/each}
            </table>
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
    
    .section table, th, td {
        border: 1px solid black;
        border-collapse: collapse;
    }
    .section table {
        width: 100%;
    }
    
    .section th, td {
        text-align: center;
        padding: 5px;
    }
    

    .workers a {
        color: black;
    }
</style>
