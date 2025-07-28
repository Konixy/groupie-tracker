<script lang="ts">
	let { firstRender }: { firstRender: boolean } = $props();

	$inspect(firstRender);
</script>

{#if firstRender}
	<div class="loading">Chargement...</div>
{:else}
	<div class="logo">Groupie Tracker</div>
{/if}

<style>
	.logo {
		font-family: 'Pacifico', cursive;
		font-size: 5rem;
		margin-bottom: 0;
		line-height: 1;
		text-align: center;
	}

	.loading {
		font-family: 'Pacifico', cursive;
		font-size: 5rem;
		margin-bottom: 0;
		line-height: 1;
		text-align: center;
		color: var(--light-muted);
		animation: loading 1.5s infinite;
	}

	@keyframes loading {
		0% {
			opacity: 0.5;
		}
		50% {
			opacity: 1;
		}
		100% {
			opacity: 0.5;
		}
	}
</style>
