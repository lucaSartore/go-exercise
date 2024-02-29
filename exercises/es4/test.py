import pandas as pd
from datetime import datetime, date, timedelta
import datetime as dt
import time
data = pd.read_json("exercises/es4/dataset.json")

data["start"] = pd.to_datetime(data["start"])
data["finish"] = pd.to_datetime(data["finish"],format='mixed')

print(data)

days = list(data["start"].unique()) + list(data["finish"].unique())
days = [day.date() for day in set(days)]
days = set(days)
ids = set(data["id"])

print("Days:", days)
print("Ids:", ids)


t1 = time.time()

for id in ids:
    filtered = data[data["id"] == id]
    for day in days:
        day: date = day

        start = datetime.combine(day, dt.time.min)
        end = datetime.combine(day, dt.time.max)


        filtered_day: pd.DateOffset = filtered[
            (filtered["start"] <= end) &
            (filtered["finish"] >= start)
        ].copy()

        def get_weight(row):
            time_before = (start - row["start"]).total_seconds()
            time_after = (row["finish"] - end).total_seconds()
            
            time_before = max(time_before, 0)
            time_after = max(time_after, 0)

            total_time = (row["finish"] - row["start"]).total_seconds()

            time_inside = total_time - time_before - time_after
        


            if time_inside  == 0:
                if total_time == 0:
                    w = 1
                else:
                    w = 0
            else:
                w = time_inside / total_time

            row["weight"] = w
            return row
        
        if filtered_day.empty:
            continue

        filtered_day =  filtered_day.apply(lambda x: get_weight(x), axis=1)

        filtered_day["final_x"] = filtered_day["x"] * filtered_day["weight"]

        x = filtered_day["final_x"].sum()
        
        #print(f"Id: {id}, Day: {day}, X: {x}")

delta_t = time.time() - t1

print(f"Time: {delta_t}")
