const queueLengthEl = document.getElementById('queue-length');
const tasksDispatchedEl = document.getElementById('tasks-dispatched');
const activeNodesEl = document.getElementById('active-nodes');
const nodeGrid = document.getElementById('node-grid');

const MAX_LOAD_FOR_BAR = 5;

function loadColor(taskCount) {
  const ratio = Math.min(taskCount / MAX_LOAD_FOR_BAR, 1);
  if (ratio < 0.4) return '#22c55e';
  if (ratio < 0.75) return '#f59e0b';
  return '#ef4444';
}

function renderNodes(nodes) {
  nodeGrid.innerHTML = nodes.map(n => {
    const status = (n.Status || 'idle').toLowerCase();
    const widthPct = Math.min((n.TaskCount / MAX_LOAD_FOR_BAR) * 100, 100);
    return `
      <div class="node-card ${status}">
        <div class="node-id">${n.ID}</div>
        <div class="node-status status-${status}">${status}</div>
        <div class="load-bar-track">
          <div class="load-bar-fill" style="width:${widthPct}%; background:${loadColor(n.TaskCount)};"></div>
        </div>
        <div class="task-count">${n.TaskCount} active task(s)</div>
      </div>
    `;
  }).join('');

  const active = nodes.filter(n => (n.Status || '').toLowerCase() !== 'offline').length;
  activeNodesEl.textContent = active;
}

const evtSource = new EventSource('/api/stream');

evtSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  queueLengthEl.textContent = data.queue_length;
  tasksDispatchedEl.textContent = data.tasks_dispatched;
  renderNodes(data.nodes || []);
};

evtSource.onerror = () => {
  console.warn('Stream disconnected, browser will auto-retry');
};