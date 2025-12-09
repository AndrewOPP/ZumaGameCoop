export class Renderer {
    // Явно объявляем свойства класса с типом any
    private canvas: any;
    private ctx: any;

    /**
     * @param {any} canvasId - ID элемента canvas, или сам элемент (зависит от логики).
     * @param {any} width - Ширина canvas.
     * @param {any} height - Высота canvas.
     */
    constructor(canvasElement: any, width: any, height: any) {
        // Логика конструктора
        // this.canvas = document.getElementById(canvasId); // Если передается ID

        this.canvas = canvasElement;
        
        // --- ВАЖНОЕ ПРИМЕЧАНИЕ ---
        // Если вы передаете сюда HTMLCanvasElement (как в вашем React-коде), 
        // то эта строка вызовет ошибку, поскольку document.getElementById ждет string.
        // Если вы изменили логику на прием элемента: this.canvas = canvasId;

        this.ctx = this.canvas.getContext("2d");
        this.canvas.width = width;
        this.canvas.height = height;
    }
    
    /**
     * Отрисовывает текущее состояние игры.
     * @param {any} gameState - Объект, содержащий данные для отрисовки.
     */
    draw(gameState: any): void { // Возвращаемый тип - void
        // Приведение типов для безопасного использования (хотя мы используем any, 
        // TypeScript знает методы canvas и ctx)
        const canvas: HTMLCanvasElement = this.canvas;
        const ctx: CanvasRenderingContext2D = this.ctx;
        
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        
        // Старая логика отрисовки тестовой координаты
        const x: any = gameState.TestCoordinate;

        if (typeof x === "number") {
            const y: number = canvas.height / 2;
            // Тип gameState.CurrentBall.r не определен, оставляем any для ballRadius
            const ballRadius: any = gameState.CurrentBall.r; 
            const drawX: number = (x % canvas.width);
            console.log(gameState);
            
            ctx.beginPath();
            ctx.arc(drawX, y, ballRadius, 0, Math.PI * 2);
            ctx.fillStyle = gameState.CurrentBall.color;
            ctx.fill();
            ctx.closePath();
        }

        // TODO: Здесь будет новая логика отрисовки Path и Chain
    }
}