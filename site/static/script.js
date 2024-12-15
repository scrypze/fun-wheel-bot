const canvas = document.getElementById('wheel');
const ctx = canvas.getContext('2d');
let items = [];
let spinning = false;

function drawWheel() {
    const radius = canvas.width / 2;
    const center = { x: radius, y: radius };
    const sliceAngle = (2 * Math.PI) / items.length;

    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.font = "14px Arial";

    items.forEach((item, i) => {
        const startAngle = i * sliceAngle;
        const endAngle = startAngle + sliceAngle;

        ctx.beginPath();
        ctx.moveTo(center.x, center.y);
        ctx.arc(center.x, center.y, radius, startAngle, endAngle);
        ctx.closePath();

        ctx.fillStyle = i % 2 === 0 ? '#FFDD57' : '#FFAB57';
        ctx.fill();
        ctx.stroke();

        ctx.save();
        ctx.translate(center.x, center.y);
        ctx.rotate(startAngle + sliceAngle / 2);
        ctx.textAlign = "right";
        ctx.fillStyle = "#000";
        ctx.fillText(item, radius - 10, 5);
        ctx.restore();
    });
}

async function loadItems() {
    const response = await fetch('/get-items');
    const data = await response.json();
    items = data.items;
    drawWheel();
}

async function addItem() {
    const input = document.getElementById('item');
    const item = input.value;
    if (item) {
        await fetch('/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text: item })
        });
        input.value = '';
        await loadItems();
    }
}

async function resetItems() {
    await fetch('/reset', { method: 'POST' });
    items = [];
    drawWheel();
    document.getElementById('result').innerText = '';
}

async function spinWheel() {
    if (spinning) return;
    if (items.length === 0) {
        alert('Добавьте элементы перед вращением!');
        return;
    }

    spinning = true;
    const response = await fetch('/spin');
    const data = await response.json();
    items = data.items;
    const winnerIndex = items.indexOf(data.winner);
    let currentAngle = 0;
    const targetAngle = (2 * Math.PI) * 5 + winnerIndex * (2 * Math.PI / items.length);

    const spin = setInterval(() => {
        currentAngle += 0.1;
        if (currentAngle >= targetAngle) {
            clearInterval(spin);
            spinning = false;
            document.getElementById('result').innerText = `Победитель: ${data.winner}`;
        }
        ctx.save();
        ctx.translate(canvas.width / 2, canvas.height / 2);
        ctx.rotate(currentAngle);
        ctx.translate(-canvas.width / 2, -canvas.height / 2);
        drawWheel();
        ctx.restore();
    }, 16);
}

async function removeLastWinner() {
    await fetch('/remove-winner', { method: 'POST' });
    await loadItems();
    document.getElementById('result').innerText = '';
}

loadItems();
