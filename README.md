# go-retries
##### This project is a library to do retries easily. You can customize configuration of delay, max retries and configuration errors unrecoverable  

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
    //Setting max retries, delay time and errors unrecoverable
    retry.Setting(map[Config]int{
    		ConfigMaxRetries: 5,
    		ConfigDelaySec: 2,
    	}, []error{errors.New("Unrecoverable")})
    
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