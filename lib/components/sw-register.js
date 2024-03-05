export let registration;

if ('serviceWorker' in navigator) {
	try {
		registration = await navigator.serviceWorker.register('/assets/js/lib/components/sw.js');
	} catch (error) {
		console.error(`Registration failed with ${error}`);
	}
}

