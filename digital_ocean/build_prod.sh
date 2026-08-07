docker stop expenseowl
docker rm expenseowl

#docker rmi expenseowl

#rm -rf ExpenseOwl

#git clone git@github.com:kwokgordon/ExpenseOwl.git

cd ExpenseOwl

git pull

docker build -t expenseowl .

cd ..

docker run -itd --restart=always -v expenseowl-data:/app/data --name expenseowl -p 127.0.0.1:8080:8080 expenseowl
