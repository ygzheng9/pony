require("expose-loader?$!expose-loader?jQuery!jquery");

(function() {
  "use strict";

  var treeviewMenu = $(".app-menu");

  // Toggle Sidebar
  $('[data-toggle="sidebar"]').click(function(event) {
    event.preventDefault();
    $(".app").toggleClass("sidenav-toggled");
  });

  // Activate sidebar treeview toggle
  $("[data-toggle='treeview']").click(function(event) {
    event.preventDefault();

    if (
      !$(this)
        .parent()
        .hasClass("is-expanded")
    ) {
      treeviewMenu
        .find("[data-toggle='treeview']")
        .parent()
        .removeClass("is-expanded");
    }
    $(this)
      .parent()
      .toggleClass("is-expanded");
  });

  // Set initial active toggle
  $("[data-toggle='treeview.'].is-expanded")
    .parent()
    .toggleClass("is-expanded");

  //Activate bootstrip tooltips
  $("[data-toggle='tooltip']").tooltip();

  // 点击第一级菜单
  $(".app-menu__item").on("click", function() {
    // 清空所有之前选中
    $(".app-menu__item").removeClass("active");

    // 选中当前
    const self = $(this);
    self.addClass("active");
  });

  // 选中当前页面对应的菜单项
  function seleteMenu() {
    // 每个页面，都返回自己的的菜单项
    const menu = $("#idData").data("server");
    console.log(menu);

    //  第一级菜单
    const items1 = $(".app-menu__item");
    items1.removeClass("active");

    let isFind = false;
    items1.each((_, i) => {
      if (i.text === menu) {
        $(i).addClass("active");
        isFind = true;
      }
    });

    //  清除之前所有的选择
    const items2 = $(".treeview-item");
    items2.removeClass("active");

    // 在第一级已经找到，直接返回
    if (isFind) {
      return;
    }

    // 第一级中没找到，再到第二级中找；
    items2.each(function() {
      const self = $(this);

      if (self.text().trim() === menu) {
        // 子菜单选中
        self.addClass("active");

        // 上一级菜单全部收起来
        $(".treeview").removeClass("is-expanded");
        // 当前选中菜单的上一级，打开
        self.parents(".treeview").addClass("is-expanded");
      }
    });
  }

  seleteMenu();
})();
