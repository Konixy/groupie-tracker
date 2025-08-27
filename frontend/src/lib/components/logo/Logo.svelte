<script lang="ts">
	import { fade } from 'svelte/transition';
	import LogoContent from './LogoContent.svelte';

	let { loading }: { loading: boolean } = $props();
</script>

{#if loading}
	<div class="loading" out:fade={{ duration: 200 }}>Chargement...</div>
{/if}

<LogoContent show={!loading} />

<style>
	.loading {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		font-family: 'Pacifico', cursive;
		font-size: 3rem;
		margin-bottom: 0;
		line-height: 1;
		text-align: center;
		color: var(--light-muted);
		animation: loading 1.5s infinite;
	}

	@keyframes loading {
		0%,
		100% {
			opacity: 0.5;
		}
		50% {
			opacity: 1;
		}
	}
</style>
