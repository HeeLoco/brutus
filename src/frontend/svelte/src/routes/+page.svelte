<script lang="ts">
	let count = 0;
	
	function increment() {
		count += 1;
	}
</script>

<h1>Azure Resource Manager</h1>
<p>SvelteKit Frontend for Multi-Backend API</p>

<button on:click={increment}>
	Count: {count}
</button>

<div class="info">
	<p>This frontend can connect to:</p>
	<ul>
		<li>Go Backend (Port 8080)</li>
		<li>Python Backend (Port 8000)</li>
		<li>TypeScript Backend (Port 3000)</li>
	</ul>
</div>

<style>
	h1 {
		color: #ff3e00;
		font-size: 2rem;
		font-weight: 100;
	}
	
	.info {
		margin-top: 2rem;
		padding: 1rem;
		background: #f8f9fa;
		border-radius: 0.5rem;
	}
	
	button {
		background: #ff3e00;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 0.25rem;
		cursor: pointer;
	}
	
	button:hover {
		background: #e63800;
	}
</style>