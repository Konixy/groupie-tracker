<script>
	// Import de la police Montserrat
	const link = document.createElement('link');
	link.rel = 'stylesheet';
	link.href =
		'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700&display=swap';
	document.head.appendChild(link);
</script>

<footer class="footer">
	<p>
		Brought to you by
		<a href="https://github.com/Ramswed" target="_blank" rel="noopener noreferrer">Manon</a>
		&
		<a href="https://github.com/Konixy" target="_blank" rel="noopener noreferrer">Anatole</a>
	</p>
</footer>

<style>
	.footer {
		padding: 1rem;
		text-align: center;
		font-family: 'Montserrat', sans-serif !important;
		font-weight: 700;
	}

	.footer p {
		font-family: 'Montserrat', sans-serif !important;
		font-weight: 700;
	}

	.footer a {
		font-family: 'Montserrat', sans-serif !important;
		font-weight: 700;
	}

	.footer a {
		color: var(--light-vibrant);
		text-decoration: none;
		transition: color 0.2s ease;
	}

	.footer a:hover {
		color: var(--vibrant);
		text-decoration: underline;
	}
</style>
