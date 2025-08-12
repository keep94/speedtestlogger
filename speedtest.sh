if speedtest -f csv &> stl_out.csv; then
    ~/go/bin/stllog -db ~/stl/stl.db -csv stl_out.csv
else
    ~/go/bin/stllog -db ~/stl/stl.db
fi
