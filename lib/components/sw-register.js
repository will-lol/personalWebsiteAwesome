export let registration;

if ('serviceWorker' in navigator) {
	try {
		registration = await navigator.serviceWorker.register('assets/js/util/components/sw.js');
	} catch (error) {
		console.error(`Registration failed with ${error}`);
	}
}

