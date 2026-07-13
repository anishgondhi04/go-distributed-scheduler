const canvas = document.getElementById('network');
const ctx = canvas.getContext('2d');

let nodePositions = [];
let particles = [];
let nodesState = [];

function layoutNodes(nodes) {
  const centerX = canvas.width / 2;
  const centerY = 60;
  const radius = 180;
  nodePositions = nodes.map((n, i) => {
    const angle = (i / nodes.length) * Math.PI * 2 - Math.PI / 2;
    return {
      id: n.ID,
      x: centerX + Math.cos(angle) * radius,
      y: 220 + Math.sin(angle) * 140,
    };
  });
}

function spawnParticle(nodeId) {
  const target = nodePositions.find(p => p.id === nodeId);
  if (!target) return;
  particles.push({
    x: canvas.width / 2,
    y: 60,
    targetX: target.x,
    targetY: target.y,
    progress: 0,
  });
}

function draw() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  ctx.fillStyle = '#38bdf8';
  ctx.beginPath();
  ctx.arc(canvas.width / 2, 60, 14, 0, Math.PI * 2);
  ctx.fill();
  ctx.font = '12px sans-serif';
  ctx.fillStyle = '#94a3b8';
  ctx.textAlign = 'center';
  ctx.fillText('scheduler', canvas.width / 2, 90);

  nodePositions.forEach(p => {
    const state = nodesState.find(n => n.ID === p.id);
    const busy = state && state.Status === 'busy';

    ctx.strokeStyle = 'rgba(148,163,184,0.15)';
    ctx.beginPath();
    ctx.moveTo(canvas.width / 2, 60);
    ctx.lineTo(p.x, p.y);
    ctx.stroke();

    ctx.beginPath();
    ctx.arc(p.x, p.y, busy ? 22 : 18, 0, Math.PI * 2);
    ctx.fillStyle = busy ? '#f59e0b' : '#22c55e';
    ctx.shadowColor = busy ? '#f59e0b' : '#22c55e';
    ctx.shadowBlur = busy ? 20 : 8;
    ctx.fill();
    ctx.shadowBlur = 0;

    ctx.fillStyle = '#e2e8f0';
    ctx.font = '11px sans-serif';
    ctx.fillText(p.id, p.x, p.y + 34);
  });

  particles.forEach(pt => {
    pt.progress += 0.03;
    const x = pt.x + (pt.targetX - pt.x) * pt.progress;
    const y = pt.y + (pt.targetY - pt.y) * pt.progress;
    ctx.beginPath();
    ctx.arc(x, y, 4, 0, Math.PI * 2);
    ctx.fillStyle = '#facc15';
    ctx.shadowColor = '#facc15';
    ctx.shadowBlur = 10;
    ctx.fill();
    ctx.shadowBlur = 0;
  });

  particles = particles.filter(pt => pt.progress < 1);
  requestAnimationFrame(draw);
}

requestAnimationFrame(draw);

const queueLengthEl = document.getElementById('queue-length');
const tasksDispatchedEl = document.getElementById('tasks-dispatched');
const activeNodesEl = document.getElementById('active-nodes');
const nodeGrid = document.getElementById('node-grid');
const strategySelect = document.getElementById('strategy-select');

const MAX_LOAD_FOR_BAR = 5;

function loadColor(taskCount) {
  const ratio = Math.min(taskCount / MAX_LOAD_FOR_BAR, 1);
  if (ratio < 0.4) return '#22c55e';
  if (ratio < 0.75) return '#f59e0b';
  return '#ef4444';
}

function renderNodes(nodes) {
  if (nodePositions.length !== nodes.length) {
    layoutNodes(nodes);
  }

  const previous = nodesState;
  nodes.forEach(n => {
    const prev = previous.find(p => p.ID === n.ID);
    if (prev && n.TaskCount > prev.TaskCount) {
      spawnParticle(n.ID);
    }
  });
  nodesState = nodes;

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

fetch('/api/strategy')
  .then(res => res.json())
  .then(data => { strategySelect.value = data.strategy; });

strategySelect.addEventListener('change', () => {
  fetch('/api/strategy', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ strategy: strategySelect.value }),
  });
});