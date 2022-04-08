<script lang="ts">
	import axios,{AxiosResponse} from 'axios';
	let key:string = "key";
	let data:string = `{"your":"data"}`;
	let tokenInput:string;
	let keyTokenInput:string;
	let tokenData:string|undefined;
	let timeData:string|undefined;
	let outputTokenReverse: any;
	let errorTokenInput:boolean = false
	let error:boolean = false

	const submitData = () => {
		type Output = {token:string}
		axios.post<Output,AxiosResponse<Output>>("/api/jwt",{key,data}).
				then(objData => {
					error = false
					tokenData = objData.data.token
		}).catch((e) => {
			console.log(e)
			error = true
		})
	}

	const reverseToken = () => {
		type Output = {body:any,time:string}
		axios.post<Output,AxiosResponse<Output>>("/api/jwt-decrypt",{key:keyTokenInput,token:tokenInput}).
				then(objData => {
					errorTokenInput = false
			outputTokenReverse = objData.data.body
			timeData = new Date(objData.data.time).toString()
		}).catch((e) => {
			console.log(e)
			errorTokenInput = true
		})
	}


</script>

<main>
	<h1>JWT Simulator with HS256 Algorithm</h1>
	<label>
		Key
		<input type="text" bind:value={key}/>
	</label>
	<label>
		Data
		<input type="text" bind:value={data}/>
	</label>
	<button on:click={submitData}>Get JWT</button>
	<h2>Output Token: </h2>
	{#if tokenData !== undefined}
		<h3 style="width: 50%;word-wrap: break-word;">Token : {tokenData}</h3>
	{/if}
	{#if error}
		<h3>Error</h3>
	{/if}
	<br>
	<h1>Reverse Token</h1>
	<label>
		Key
		<input type="text" bind:value={keyTokenInput}/>
	</label>
	<label>
		Token
		<input type="text" bind:value={tokenInput}/>
	</label>
	<button on:click={reverseToken}>Reverse</button>
	<h2>Output Reverse</h2>
	{#if outputTokenReverse !== undefined}
		<h3>Body : {outputTokenReverse}</h3>
		<h3>Time : {timeData}</h3>
	{/if}
	{#if errorTokenInput}
		<h3>Error reversing</h3>
	{/if}
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>