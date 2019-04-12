require("expose-loader?$!expose-loader?jQuery!jquery");
// require("bootstrap/dist/js/bootstrap.bundle.js");

$(() => {
  function beforeSummit(event) {
    const inputs = document.getElementsByTagName("input");

    // 原始答案
    const inputList = Array.prototype.slice.call(inputs);
    const raws = inputList
      .filter(i => i.name !== "authenticity_token")
      .map(i => ({
        name: i.name,
        value: i.value,
        checked: i.checked,
        type: i.type
      }));

    // 按照题目，归集答案
    const answers = raws.reduce((acc, curr) => {
      const idx = acc.findIndex(a => a.name === curr.name);
      if (idx === -1) {
        if (curr.type === "checkbox" || curr.type === "radio") {
          // 选中的记录答案，没选中的，不需要记录
          if (curr.checked) {
            acc.push({ name: curr.name, values: [curr.value] });
          } else {
            acc.push({ name: curr.name, values: [] });
          }
        } else {
          // 文本框
          const v = curr.value.trim();
          if (v.length > 0) {
            acc.push({ name: curr.name, values: [v] });
          } else {
            acc.push({ name: curr.name, values: [] });
          }
        }

        return acc;
      }

      if (curr.type === "checkbox" || curr.type === "radio") {
        if (curr.checked) {
          acc[idx].values.push(curr.value);
        }
      } else {
        // 文本框
        const v = curr.value.trim();
        if (v.length > 0) {
          acc[idx].values.push(v);
        }
      }

      return acc;
    }, []);

    const missings = answers
      .filter(i => i.values.length === 0)
      .map(i => i.name.toString());
    if (missings.length > 0) {
      const msg = `请回答下列问题：${missings.join(" ")}`;
      layer.alert(msg);

      $("#errorInfo").text(msg);
      $("#alertInfo").css("display", "block");

      // 把没回答的问题高亮
      const total = $(".question-title");
      total.removeClass("alert alert-danger");

      missings.forEach(s => {
        total.each(function() {
          const title = $(this);
          const i = title.data("seq").toString();
          if (i === s) {
            title.addClass("alert alert-danger");
          }
        });
      });

      event.preventDefault();
      return false;
    }

    return true;
  }
  const btn = document.getElementById("btnSubmit");
  btn.addEventListener("click", beforeSummit);

  // 选项排序
  function moveUpDown(btn, up) {
    const option = $(btn.parents(".option")[0]);

    if (up) {
      option.insertBefore(option.prev());
    } else {
      option.insertAfter(option.next());
    }

    // 获取选项的顺序
    const options = [];
    const question = $(option.parents(".sort")[0]);
    question.children(".option").each(function() {
      const opt = $(this);
      options.push(opt.data("seq"));
    });
    const inputID = question.data("target");
    const inputBox = document.getElementById(inputID);
    inputBox.value = options.join(" ");
  }

  // 向上的按钮
  $(".action-up").on("click", function() {
    moveUpDown($(this), true);
  });

  // 向下的按钮
  $(".action-down").on("click", function() {
    moveUpDown($(this), false);
  });
});
