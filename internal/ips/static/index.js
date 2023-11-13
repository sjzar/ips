let editor;

document.addEventListener("DOMContentLoaded", function() {
    const button = document.querySelector(".confirm-button");
    const inputBox = document.querySelector(".input-box");

    editor = CodeMirror.fromTextArea(document.getElementById("jsonEditor"), {
        mode: "application/ld+json",
        theme: "material",
        lineNumbers: true,
        readOnly: true
    });

    let timer;

    function queryIP() {
        clearTimeout(timer);

        const inputText = encodeURI(inputBox.value);

        if (inputText) {
            fetch(`/api/v1/query?text=${inputText}`)
                .then(response => response.json())
                .then(data => {
                    editor.setValue(JSON.stringify(data, null, 4));
                })
                .catch(error => {
                    editor.setValue("查询过程中发生错误，请稍后重试。");
                });
        } else {
            editor.setValue("请输入有效信息！");
        }
    }

    button.addEventListener("click", queryIP);

    inputBox.addEventListener("input", function() {
        clearTimeout(timer); // 清除之前的定时器
        timer = setTimeout(queryIP, 1000); // 1 秒后触发查询
    });

    inputBox.addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            queryIP();
        }
    });

    inputBox.addEventListener("paste", function(event) {
        setTimeout(function() {
            queryIP();
        }, 0);
    });

    const bodyElement = document.body;
    const toggleIcon = document.getElementById("toggleIcon");
    toggleIcon.addEventListener("click", function() {
        if (bodyElement.getAttribute("data-theme") === "dark") {
            bodyElement.removeAttribute("data-theme");
            toggleIcon.className = "fa fa-moon-o";
        } else {
            bodyElement.setAttribute("data-theme", "dark");
            toggleIcon.className = "fa fa-sun-o";
        }
    });

});
