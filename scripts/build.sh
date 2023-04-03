# * Удаление папки: api, out, ./main
rm -rf apis/
rm -rf out
rm ./main

# * Создаёт папку apis 
mkdir apis

echo "Compilation core os 🔥"
# * Компиляция main.go
go build ./cmd/main.go

echo "Launching core os 🔥"
# * Запусит ядра
./main

# * tidy удостоверяется, что go. mod соответствует исходному коду в модуле. 
# * Он добавляет все недостающие модули, необходимые для построения пакетов 
# * и зависимостей текущего модуля, и удаляет неиспользуемые модули, которые 
# * не предоставляют никаких соответствующих пакетов. Он также добавляет все
# * недостающие записи в go
go mod tidy

echo ""
echo "Compilation server os 🔥"
# * Компилирует и запускает server os
go build  -o out/app out/app.go

echo "Launching server os 🔥"
./out/app