docker stop uat_expenseowl
docker rm uat_expenseowl

#docker rmi uat_expenseowl

#rm -rf ExpenseOwl

#git clone git@github.com:kwokgordon/ExpenseOwl.git

cd ExpenseOwl

git switch uat
git pull

docker build -t uat_expenseowl .

cd ..

docker run -itd --restart=always -v uat-expenseowl-data:/app/data --name uat_expenseowl -p 3000:8080 uat_expenseowl
