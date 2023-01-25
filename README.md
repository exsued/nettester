Робот, регулярно проверяющий наличие интернета (Raspberry Pi)
--------------------------------------

Инструкция по установке:
------------------------

#### 1. Установка Raspbian:
	1. Установить rpi-imager: sudo apt install rpi-imager
	2. Запустить rpi-imager
	3. Выбрать CHOOSE OS -> Raspberry Pi OS (other) -> Raspberry Pi OS Lite (32-bit)
	4. Вставить SD карту в компьютер и нажать write
	5. По завершении установки вставить SD карту в Raspberry Pi
#### 2. Установка net tester:
	1. Зайти в систему raspbian и подключить её к интернету
	2. Выполнить: wget https://github.com/exsued/nettester/releases/download/v1.00-alpha/nettester_install.py
	3. Выполнить: chmod +x ./nettester_install.py
	4. Выполнить: ./nettester_install.py
	5. Enter device name: <ввести имя устройства для tcp сервера>
