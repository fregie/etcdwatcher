# etcdwatcher
etcdwatcher is a simple library to watch etcd changes.Usually used to monitor configuration items.
Just define the key to watch and the value type,it will update the value.

## Installation
```bash
go get github.com/fregie/etcdwatcher
```

## usage
```go
import "github.com/fregie/etcdwatcher"
```
watch a simple configuration:
```go
// Create a new watcher
etcdEndpoint := []string{"http://localhost:2379"}
configWatcher, err := etcdwatcher.NewWatcher(etcdEndpoint)
if err != nil {
  panic()
}
// Define the key to watch
var (
  int32Value *etcdwatcher.Int32 = etcdwatcher.NewInt32("/etcdwatcher/test/int32Value", 0)
  stringValue *etcdwatcher.String = etcdwatcher.NewString("/etcdwatcher/test/int32String", "default")
  durationValue *etcdwatcher.Duration = etcdwatcher.NewDuration("/etcdwatcher/test/durationValue", time.Second)
  boolValue *etcdwatcher.Bool = etcdwatcher.NewBool("/etcdwatcher/test/boolValue", false)
)

// Watch the key
err = configWatcher.WatchItems([]etcdwatcher.Item{
  int32Value,
  stringValue,
  durationValue,
  boolValue,
})
if err != nil {
  panic()
}

// Use
fmt.Println(int32Value.Value())
fmt.Println(stringValue.Value())
fmt.Println(durationValue.Value())
fmt.Println(boolValue.Value())
```

## custom your watch item
etcdwatcher watch a item implement the Item interface.  
You can implement your own Item to watch your custom item.
```go
type Item interface {
  // Return the key watching
  Key() string
  // Parse the data when the value is changed
  Parse([]byte) error
  // Set the default value
  SetDefault()
}
```
### example
define a `map[string]string` Item to watch
```go
type customMapItem struct {
  sync.Map
  key string
}

func NewcustomMapItem(key string) *customMapItem {
  return &customMapItem{
    key: key,
  }
}

func (c *customMapItem) SetDefault() {
  c.Range(func(key, value interface{}) bool {
    c.Delete(key)
    return true
  })
}

func (c *customMapItem) Parse(data []byte) error {
  dataMap := make(map[string]string)
  err := json.Unmarshal(data, &dataMap)
  if err != nil {
    return err
  }
  for k, v := range dataMap {
    c.Store(k, v)
  }
  return nil
}

func (c *customMapItem) Key() string {
  return c.key
}

func (c *customMapItem) SaveToEtcd(ctx context.Context, cli *clientv3.Client, dataMap map[string]string) error {
  data, err := json.Marshal(dataMap)
  if err != nil {
    return err
  }
  _, err := cli.Put(ctx, m.Key(), string(data))
  return err
}
```