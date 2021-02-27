# go-retries
##### This project is a library to do retries easily. You can customize configuration of delay, max retries and configuration errors unrecoverable  

The default settings
*   Max retries: 3
*   Delay seconds: 3

Simple use example below:
```go
import retry "github.com/arturmartini/go-retries"

func example(any string) error {
    errRetry := retry.New().Do(func() error) error {
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
    errRetry := New().SetConfigurations(
    		Configuration{Key: ConfigMaxRetries, Value: 5},
    		Configuration{Key: ConfigDelaySec, Value: 5}).
            //Setting recoverable errors
    		SetRecoverableErrors(errors.New("any error")).Do(func() interface{} {
    		var response, err = http.Get("https://example.com")
    		doSomething(response)
    		return err		
    	})
    
    if errRetry != nil {
        //do something
        return errRetry
    }
    return nil   
}
```