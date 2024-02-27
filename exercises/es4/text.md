# Exercise text

Given the data in dataset.json
```json
 {"start":"2022-09-05 15:04:43.195098","finish":"2022-09-06 17:04:45.195098","id":"1","x":"5"},
```

you should find the value of x for every day and id, notin that
 - the same day and ip can be present on two different rows. in this case you should sum the values
 - start and finish are not always in the same day. if this is the case x need to be split in multiple days proportionally  