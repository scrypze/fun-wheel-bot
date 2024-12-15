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

async function addItem() {
    const item = document.getElementById('item').value;
    if (item) {
        await fetch('/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text: item })
        });
        items.push(item);
        drawWheel();
    }
}

async function resetItems() {
    await fetch('/reset', { method: 'POST' });
    items = [];
    drawWheel();
}

async function spinWheel() {
    if (spinning) return;

    spinning = true;
    const response = await fetch('/spin');
    const data = await response.json();
    const winnerIndex = items.indexOf(data.winner);
    let currentAngle = 0;
    const targetAngle = (2 * Math.PI) * 5 + winnerIndex * (2 * Math.PI / items.length);

    const spin = setInterval(() => {
        currentAngle += 0.1;
        if (currentAngle >= targetAngle) {
            clearInterval(spin);
            spinning = false;
            document.getElementById('result').innerText = `Winner: ${data.winner}`;
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
    items = items.filter(item => item !== document.getElementById('result').innerText.split(': ')[1]);
    drawWheel();
    document.getElementById('result').innerText = '';
}

drawWheel();
