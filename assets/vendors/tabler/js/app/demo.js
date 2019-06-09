// 定义 module，返回一个 obj 即可
define([], function () {
    function func1(msg) {
        console.log("demo for test1: ", msg);
    }

    return {
        func1
    }
});