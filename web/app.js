const app = document.getElementById('app');
app.innerHTML = '<div id="live"></div>';
const live = document.getElementById('live');

const evtSource = new EventSource('/api/stream');

evtSource.onmessage = (event) => {
	const data = JSON.parse(event.data);
	live.innerHTML = `
		<p>Queue length: ${data.queue_length}</p>
		<p>Tasks dispatched: ${data.tasks_dispatched}</p>
		<ul>
			${data.nodes.map(n => `<li>${n.ID}: ${n.Status} (${n.TaskCount} tasks)</li>`).join('')}
		</ul>
	`;
};

evtSource.onerror = () => {
	console.warn('Stream disconnected, browser will auto-retry');
};