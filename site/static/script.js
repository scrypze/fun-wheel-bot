const canvas = document.getElementById('wheel');
const ctx = canvas.getContext('2d');
let items = [];
let spinning = false;

function drawWheel() {
    if (items.length === 0) {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        return;
    }

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
    try {
        const response = await fetch('/items');
        if (!response.ok) {
            throw new Error('Ошибка загрузки элементов');
        }
        const data = await response.json();
        items = data.items;
        drawWheel();
    } catch (error) {
        console.error('Ошибка при загрузке элементов:', error);
    }
}

async function addItem() {
    const input = document.getElementById('item');
    const item = input.value.trim();
    if (item) {
        try {
            const response = await fetch('/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ text: item })
            });
            if (!response.ok) {
                throw new Error('Ошибка добавления элемента');
            }
            input.value = '';
            await loadItems();
        } catch (error) {
            console.error('Ошибка при добавлении элемента:', error);
            alert('Не удалось добавить элемент');
        }
    }
}

async function resetItems() {
    try {
        const response = await fetch('/reset', { method: 'POST' });
        if (!response.ok) {
            throw new Error('Ошибка сброса элементов');
        }
        items = [];
        drawWheel();
        document.getElementById('result').innerText = '';
    } catch (error) {
        console.error('Ошибка при сбросе элементов:', error);
        alert('Не удалось сбросить элементы');
    }
}

async function spinWheel() {
    if (spinning) return;
    if (items.length === 0) {
        alert('Добавьте элементы перед вращением!');
        return;
    }

    try {
        spinning = true;
        const response = await fetch('/spin');
        if (!response.ok) {
            throw new Error('Ошибка вращения колеса');
        }
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
    } catch (error) {
        console.error('Ошибка при вращении колеса:', error);
        alert('Не удалось провести вращение');
        spinning = false;
    }
}

async function removeLastWinner() {
    try {
        const response = await fetch('/remove-winner', { method: 'POST' });
        if (!response.ok) {
            throw new Error('Ошибка удаления победителя');
        }
        await loadItems();
        document.getElementById('result').innerText = '';
    } catch (error) {
        console.error('Ошибка при удалении победителя:', error);
        alert('Не удалось удалить победителя');
    }
}

loadItems();
