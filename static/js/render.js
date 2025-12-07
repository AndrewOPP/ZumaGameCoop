export class Renderer {
    constructor(canvasId, width, height) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext("2d");
        this.canvas.width = width;
        this.canvas.height = height;
    }
    

    draw(gameState) {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Старая логика отрисовки тестовой координаты
        const x = gameState.TestCoordinate;
        if (typeof x === "number") {
            const y = this.canvas.height / 2;
            const ballRadius = gameState.CurrentBall.r;
            const drawX = (x % this.canvas.width);
            console.log(gameState);
            
            this.ctx.beginPath();
            this.ctx.arc(drawX, y, ballRadius, 0, Math.PI * 2);
            this.ctx.fillStyle = gameState.CurrentBall.color;
            this.ctx.fill();
            this.ctx.closePath();
        }

        // TODO: Здесь будет новая логика отрисовки Path и Chain
    }
}


