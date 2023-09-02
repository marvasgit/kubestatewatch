rm -rf ./tmp
mkdir ./tmp
cd ./tmp
helm package ../diffwatcher
mkdir diffwatcher
mv *.tgz ./diffwatcher
curl https://robusta-charts.storage.googleapis.com/index.yaml > index.yaml
helm repo index --merge index.yaml --url https://robusta-charts.storage.googleapis.com ./diffwatcher
gsutil rsync -r diffwatcher gs://robusta-charts
gsutil setmeta -h "Cache-Control:max-age=0" gs://robusta-charts/index.yaml
cd ../
rm -rf ./tmp