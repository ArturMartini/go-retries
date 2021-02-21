# go-retries
##### This project is a library to do retries easily. You can customize configuration of delay, max retries and configuration errors unrecoverable  

The default settings
*   Max retries: 3
*   Delay seconds: 3

Simple use example below:
```go
import retry "github.com/arturmartini/go-retries"

func example(any string) error {
    errRetry := retry.Do(func() error) error {
        var response, err = http.Get("https://example.com")
        return err 
    }
    
    if errRetry != nil {
        //do something
        return errRetry
    }
    return nil   
}
```

Use custom configuration:
```go
import retry "github.com/arturmartini/go-retries"

func example(any string) error {
    //Setting max retries, delay time 
    retry.SetConfigurations(
    		retry.Configuration{Key: retry.ConfigMaxRetries, Value: 5},
    		retry.Configuration{Key: retry.ConfigDelaySec, Value: 5})
    
    //Setting recoverable errors
    retry.SetRecoverableErrors(errGetDataSchemaRegistry)
    
    errRetry := retry.Do(func() error) error {
        var response, err = http.Get("https://example.com")
        return err 
    }
    
    if errRetry != nil {
        //do something
        return errRetry
    }
    return nil   
}
```