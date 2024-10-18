from PIL import Image
from pystray import Icon, Menu, MenuItem

def generate_icon() -> Icon:
    image = Image.open("icons/kubernetes.png")
    icon = Icon(name="kube-context", icon=image, title="kube-context")

    return icon

def generate_icon_menu():
    context_items = []

    for i in range(0, 5):
        context_item = MenuItem(
            text="hola",
            action=on_click_select,
            checked=am_i_selected,
        )

        context_items.append(context_item)

    quit_item = MenuItem(text="Quit", action=on_click_quit)

    return Menu(
        *context_items,
        Menu.SEPARATOR,
        quit_item,
    )

def am_i_selected(menu_item: MenuItem) -> bool:
    #print(menu_item.text)
    #print(menu_item.checked)

    return True

def on_click_select():
    pass

def on_click_quit():
    global icon

    icon.stop()

icon = generate_icon()
icon.menu = generate_icon_menu()

def setup(icon):
    icon.visible = True

icon.run(setup)
