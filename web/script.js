document.addEventListener("DOMContentLoaded", function () {
    const getMenuButton = document.getElementById("getMenuButton");
    const mainCourse = document.getElementById("mainCourse");
    const sideDish = document.getElementById("sideDish");
    const fireworksContainer = document.getElementById("fireworksContainer");
    const currentTimeElement = document.getElementById("currentTime"); // 获取显示时间的元素

    getMenuButton.addEventListener("click", function () {
        getMenuButton.disabled = true; // 禁用按钮以避免多次点击

        fetch("/getMenu")
            .then(response => response.json())
            .then(data => {
                // 设置主菜和副菜
                mainCourse.textContent = "主菜：" + data.main_course;
                sideDish.textContent = "副菜：" + data.side_dish;

                // 添加动画效果
                mainCourse.classList.add("fadeIn");
                sideDish.classList.add("fadeIn");

                // 启用按钮以允许下一次选择
                getMenuButton.disabled = false;

                // 添加多个不同颜色的烟花动画
                createFireworks();

                // 显示当前时间
                showCurrentTime();
            })
            .catch(error => console.error(error));
    });

    function createFireworks() {
        const colors = ["#f39c12", "#3498db", "#e74c3c", "#2ecc71", "#9b59b6", "#f1c40f"];
        const fireworks = document.createElement("div");
        fireworks.classList.add("firework");
        fireworks.style.left = `${Math.random() * 100}%`;
        fireworks.style.top = `${Math.random() * 100}%`;
        fireworks.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)];
        fireworksContainer.appendChild(fireworks);

        setTimeout(() => {
            fireworks.remove();
        }, 3000);
    }

    // 显示当前时间的函数
    function showCurrentTime() {
        const now = new Date();
        const hours = now.getHours().toString().padStart(2, "0");
        const minutes = now.getMinutes().toString().padStart(2, "0");
        const seconds = now.getSeconds().toString().padStart(2, "0");
        const currentTime = `${hours}:${minutes}:${seconds}`;
        currentTimeElement.textContent = "当前时间：" + currentTime;

        // 每秒更新时间
        setTimeout(showCurrentTime, 1000);
    }

    // 初始化页面时显示当前时间
    showCurrentTime();
});
