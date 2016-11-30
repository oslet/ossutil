package lib

import (
	"fmt"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
	"strings"
	"time"
)

var headerOptionMap = map[string]interface{}{
	oss.HTTPHeaderContentType:             oss.ContentType,
	oss.HTTPHeaderCacheControl:            oss.CacheControl,
	oss.HTTPHeaderContentDisposition:      oss.ContentDisposition,
	oss.HTTPHeaderContentEncoding:         oss.ContentEncoding,
	oss.HTTPHeaderExpires:                 oss.Expires,
	oss.HTTPHeaderAcceptEncoding:          oss.AcceptEncoding,
	oss.HTTPHeaderOssServerSideEncryption: oss.ServerSideEncryption,
	oss.HTTPHeaderOssObjectACL:            oss.ObjectACL,
	oss.HTTPHeaderOrigin:                  oss.Origin,
}

func formatHeaderString(sep string) string {
	str := ""
	for header := range headerOptionMap {
        if header == oss.HTTPHeaderExpires {
            str += header + fmt.Sprintf("(time.RFC3339: %s)", time.RFC3339) + sep
        } else {
		    str += header + sep
        }
	}
	if len(str) >= len(sep) {
		str = str[:len(str)-len(sep)]
	}
	return str
}

func fetchHeaderOptionMap(name string) (interface{}, error) {
	for header, f := range headerOptionMap {
		if strings.ToLower(name) == strings.ToLower(header) {
			return f, nil
		}
	}
	return nil, fmt.Errorf("unsupported header: %s, please check", name)
}

func getOSSOption(name string, param string) (oss.Option, error) {
	if f, err := fetchHeaderOptionMap(name); err == nil {
		switch f.(type) {
		case func(string) oss.Option:
			return f.(func(string) oss.Option)(param), nil
		case func(oss.ACLType) oss.Option:
			return f.(func(oss.ACLType) oss.Option)(oss.ACLType(param)), nil
		case func(t time.Time) oss.Option:
            val, err := time.Parse(http.TimeFormat, param)
            if err != nil {
                val, err = time.Parse(time.RFC3339, param)
                if err != nil {
                    return nil, err
                }
            }
			return f.(func(time.Time) oss.Option)(val), nil
        default:
            return nil, fmt.Errorf("error option type, internal error")
        }
	}
	return nil, fmt.Errorf("unsupported header: %s, please check", name)
}

var specChineseSetMeta = SpecText{

	synopsisText: "设置已上传的objects的元信息",

	paramText: "url [meta] [options]",

	syntaxText: ` 
    ossutil set-meta oss://bucket[/prefix] [header:value#header:value...] [--update] [--delete] [-r] [-f] [-c file] 
`,

	detailHelpText: ` 
    该命令可设置或者更新或者删除指定objects的meta信息。当指定--recursive选项时，ossutil
    获取所有与指定url匹配的objects，批量设置这些objects的meta，否则，设置指定的单个object
    的元信息，如果该object不存在，ossutil会报错。

    （1）设置全量值：如果用户未指定--update选项和--delete选项，ossutil会设置指定objects的
        meta为用户输入的[header:value#header:value...]。当缺失[header:value#header:value...]
        信息时，相当于删除全部meta信息（对于不可删除的headers，即：不以` + oss.HTTPHeaderOssMetaPrefix + `开头的headers，
        其值不会改变）。此时ossutil会进入交互模式并要求用户确认meta信息。

    （2）更新meta：如果用户设置--update选项，ossutil会更新指定objects的指定header为输入
        的value值，其中value可以为空，指定objects的其他meta信息不会改变。此时不支持--delete
        选项。

    （3）删除meta：如果用户设置--delete选项，ossutil会删除指定objects的指定header（对于不可
        删除的headers，即：不以` + oss.HTTPHeaderOssMetaPrefix + `开头的headers，该选项不起作用），该此时value必须
        为空（header:或者header），指定objects的其他meta信息不会改变。此时不支持--update选项。

    该命令不支持bucket的meta设置，需要设置bucket的meta信息，请使用bucket相关操作。
    查看bucket或者object的meta信息，请使用stat命令。

Headers:

    可选的header列表如下：
        ` + formatHeaderString("\n        ") + `
        以及以` + oss.HTTPHeaderOssMetaPrefix + `开头的header

    注意：header不区分大小写，但value区分大小写。

用法：

    该命令有两种用法：

    1) ossutil set-meta oss://bucket/object [header:value#header:value...] [--update] [--delete] [-f] 
        如果未指定--recursive选项，ossutil设置指定的单个object的meta信息，此时请确保url
    精确指定了想要设置meta的object，当object不存在时会报错。如果指定了--force选项，则不
    会进行询问提示。如果用户未输入[header:value#header:value...]，相当于删除object的所有
    meta。
        --update选项和--delete选项的用法参考上文。

    2) ossutil set-meta oss://bucket[/prefix] [header:value#header:value...] -r [--update] [--delete] [-f]
        如果指定了--recursive选项，ossutil会查找所有前缀匹配url的objects，批量设置这些
    objects的meta信息，当错误出现时终止命令。如果--force选项被指定，则不会进行询问提示。
        --update选项和--delete选项的用法参考上文。
`,

	sampleText: ` 
    (1)ossutil set-meta oss://bucket1/obj1 Cache-Control:no-cache#Content-Encoding:gzip#X-Oss-Meta-a:b
        设置obj1的Cache-Control，Content-Encoding和X-Oss-Meta-a头域

    (2)ossutil set-meta oss://bucket1/o X-Oss-Meta-empty:#Content-Type:plain/text --update -r
        批量更新以o开头的objects的X-Oss-Meta-empty和Content-Type头域

    (3)ossutil set-meta oss://bucket1/obj1 X-Oss-Meta-delete --delete
        删除obj1的X-Oss-Meta-delete头域

    (4)ossutil set-meta oss://bucket/o -r
        批量设置以o开头的objects的meta为空
`,
}

var specEnglishSetMeta = SpecText{

	synopsisText: "set metadata on already uploaded objects",

	paramText: "url [meta] [options]",

	syntaxText: ` 
    ossutil set-meta oss://bucket[/prefix] [header:value#header:value...] [--update] [--delete] [-r] [-f] [-c file] 
`,

	detailHelpText: ` 
    The command can be used to set, update or delete the specified objects' meta data. 
    If --recursive option is specified, ossutil find all matching objects and batch set 
    meta on these objects, else, ossutil set meta on single object, if the object not 
    exist, error happens. 

    (1) Set full meta: If --update option and --delete option is not specified, ossutil 
        will set the meta of the specified objects to [header:value#header:value...], what
        user inputs. If [header:value#header:value...] is missing, it means clear the meta 
        data of the specified objects(to those headers which can not be deleted, that is, 
        the headers do not start with: ` + oss.HTTPHeaderOssMetaPrefix + `, the value will not be changed), at the 
        time ossutil will ask user to confirm the input.

    (2) Update meta: If --update option is specified, ossutil will update the specified 
        headers of objects to the values that user inputs(the values can be empty), other 
        meta data of the specified objects will not be changed. --delete option is not 
        supported in the usage. 

    (3) Delete meta: If --delete option is specified, ossutil will delete the specified 
        headers of objects that user inputs(to those headers which can not be deleted, 
        that is, the headers do not start with: ` + oss.HTTPHeaderOssMetaPrefix + `, the value will not be changed), 
        in this usage the value must be empty(like header: or header), other meta data 
        of the specified objects will not be changed. --update option is not supported 
        in the usage.

    The meta data of bucket can not be setted by the command, please use other commands. 
    User can use stat command to check the meta information of bucket or objects.

Headers:

    ossutil supports following headers:
        ` + formatHeaderString("\n        ") + `
        and headers starts with: ` + oss.HTTPHeaderOssMetaPrefix + `

    Warning: headers are case-insensitive, but value are case-sensitive.

Usage:

    There are two usages:

    1) ossutil set-meta oss://bucket/object [header:value#header:value...] [--update] [--delete] [-f] 
        If --recursive option is not specified, ossutil set meta on the specified single 
    object. In the usage, please make sure url exactly specified the object you want to 
    set meta on, if object not exist, error occurs. If --force option is specified, ossutil 
    will not show prompt question. 
        The usage of --update option and --delete option is showed in detailHelpText. 

    2) ossutil set-meta oss://bucket[/prefix] [header:value#header:value...] -r [--update] [--delete] [-f]
        If --recursive option is specified, ossutil will search for prefix-matching objects 
    and set meta on these objects, if error occurs, the operation is terminated. If --force 
    option is specified, ossutil will not show prompt question.
        The usage of --update option and --delete option is showed in detailHelpText.
`,

	sampleText: ` 
    (1)ossutil set-meta oss://bucket1/obj1 Cache-Control:no-cache#Content-Encoding:gzip#X-Oss-Meta-a:b
        Set Cache-Control, Content-Encoding and X-Oss-Meta-a header for obj1

    (2)ossutil set-meta oss://bucket1/o X-Oss-Meta-empty:#Content-Type:plain/text -u -r
        Batch update X-Oss-Meta-empty and Content-Type header on objects that start with o

    (3)ossutil set-meta oss://bucket1/obj1 X-Oss-Meta-delete -d
        Delete X-Oss-Meta-delete header of obj1 

    (4)ossutil set-meta oss://bucket/o -r
        Batch set the meta of objects that start with o to empty
`,
}

// SetMetaCommand is the command set meta for object
type SetMetaCommand struct {
	command Command
}

var setMetaCommand = SetMetaCommand{
	command: Command{
		name:        "set-meta",
		nameAlias:   []string{"setmeta", "set_meta"},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseSetMeta,
		specEnglish: specEnglishSetMeta,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionRecursion,
			OptionUpdate,
			OptionDelete,
			OptionForce,
			OptionConfigFile,
            OptionEndpoint,
            OptionAccessKeyID,
            OptionAccessKeySecret,
            OptionSTSToken,
			OptionRetryTimes,
			OptionRoutines,
            OptionLanguage,
		},
	},
}

// function for FormatHelper interface
func (sc *SetMetaCommand) formatHelpForWhole() string {
	return sc.command.formatHelpForWhole()
}

func (sc *SetMetaCommand) formatIndependHelp() string {
	return sc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (sc *SetMetaCommand) Init(args []string, options OptionMapType) error {
	return sc.command.Init(args, options, sc)
}

// RunCommand simulate inheritance, and polymorphism
func (sc *SetMetaCommand) RunCommand() error {
	isUpdate, _ := GetBool(OptionUpdate, sc.command.options)
	isDelete, _ := GetBool(OptionDelete, sc.command.options)
	recursive, _ := GetBool(OptionRecursion, sc.command.options)
	force, _ := GetBool(OptionForce, sc.command.options)
	routines, _ := GetInt(OptionRoutines, sc.command.options)
    language, _ := GetString(OptionLanguage, sc.command.options)
    language = strings.ToLower(language)

    if err := sc.checkOption(isUpdate, isDelete, force, language); err != nil {
        return err
    }

	cloudURL, err := CloudURLFromString(sc.command.args[0])
	if err != nil {
		return err
	}

	if cloudURL.bucket == "" {
		return fmt.Errorf("invalid cloud url: %s, miss bucket", sc.command.args[0])
	}

	bucket, err := sc.command.ossBucket(cloudURL.bucket)
	if err != nil {
		return err
	}

	str, err := sc.getMetaData(force, language)
	if err != nil {
		return err
	}

	headers, err := sc.parseHeaders(str, isDelete)
	if err != nil {
		return err
	}

	if !recursive {
		return sc.setObjectMeta(bucket, cloudURL.object, headers, isUpdate, isDelete)
	}
	return sc.batchSetObjectMeta(bucket, cloudURL, headers, isUpdate, isDelete, force, routines)
}

func (sc *SetMetaCommand) checkOption(isUpdate, isDelete, force bool, language string) (error) {
	if isUpdate && isDelete {
		return fmt.Errorf("--update option and --delete option are not supported for %s at the same time, please check", sc.command.args[0])
	}
    if !isUpdate && !isDelete && !force {
        if language == LEnglishLanguage {
            fmt.Printf("Warning: --update option means update the specified header, --delete option means delete the specified header, miss both options means update the whole meta info, continue to update the whole meta info(y or N)? ")
        } else {
            fmt.Printf("警告：--update选项更新指定的header，--delete选项删除指定的header，两者同时缺失会更改object的全量meta信息，请确认是否要更改全量meta信息(y or N)? ")
        }
        var str string
        if _, err := fmt.Scanln(&str); err != nil || (strings.ToLower(str) != "yes" && strings.ToLower(str) != "y") {
            return fmt.Errorf("operation is canceled")
        }
        fmt.Println("")
    }
    return nil
}

func (sc *SetMetaCommand) getMetaData(force bool, language string) (string, error) {
	if len(sc.command.args) > 1 {
		return strings.TrimSpace(sc.command.args[1]), nil
	}

	if force {
		return "", nil
	}

    if language == LEnglishLanguage {
	    fmt.Printf("Do you really mean the empty meta(or forget to input header:value pair)? \nEnter yes(y) to continue with empty meta, enter no(n) to show supported headers, other inputs will cancel operation: ")
    } else {
	    fmt.Printf("你是否确定你想设置的meta信息为空（或者忘记了输入header:value对）? \n输入yes(y)使用空meta继续设置，输入no(n)来展示支持的headers，其他输入将取消操作：")
    }
	var str string
	if _, err := fmt.Scanln(&str); err != nil || (strings.ToLower(str) != "yes" && strings.ToLower(str) != "y" && strings.ToLower(str) != "no" && strings.ToLower(str) != "n") {
		return "", fmt.Errorf("unknown input, operation is canceled")
	}
	if strings.ToLower(str) == "yes" || strings.ToLower(str) == "y" {
		return "", nil
	}

    if language == LEnglishLanguage {
	    fmt.Printf("\nSupported headers:\n    %s\n    And the headers start with: \"%s\"\n\nPlease enter the header:value#header:value... pair you want to set: ", formatHeaderString("\n    "), oss.HTTPHeaderOssMetaPrefix)
    } else {
        fmt.Printf("\n支持的headers:\n    %s\n    以及以\"%s\"开头的headers\n\n请输入你想设置的header:value#header:value...：", formatHeaderString("\n    "), oss.HTTPHeaderOssMetaPrefix)
    }
	if _, err := fmt.Scanln(&str); err != nil {
		return "", fmt.Errorf("meta empty, please check, operation is canceled")
	}
	return strings.TrimSpace(str), nil
}

func (sc *SetMetaCommand) parseHeaders(str string, isDelete bool) (map[string]string, error) {
	if str == "" {
		return nil, nil
	}

	headers := map[string]string{}
	sli := strings.Split(str, "#")
	for _, s := range sli {
		pair := strings.SplitN(s, ":", 2)
		name := pair[0]
		value := ""
		if len(pair) > 1 {
			value = pair[1]
		}
		if isDelete && value != "" {
			return nil, fmt.Errorf("delete meta for object do no support value for header:%s, please set value:%s to empty", name, value)
		}
		if _, err := fetchHeaderOptionMap(name); err != nil && !strings.HasPrefix(strings.ToLower(name), strings.ToLower(oss.HTTPHeaderOssMetaPrefix)) {
			return nil, fmt.Errorf("unsupported header:%s, please try \"help %s\" to see supported headers", name, sc.command.name)
		}
		headers[name] = value
	}
	return headers, nil
}

func (sc *SetMetaCommand) setObjectMeta(bucket *oss.Bucket, object string, headers map[string]string, isUpdate, isDelete bool) error {
	if object == "" {
		return fmt.Errorf("set object meta invalid url: %s, object empty. Set bucket meta is not supported, if you mean batch set meta on objects, please use --recursive", sc.command.args[0])
	}

	allheaders := headers
	if isUpdate || isDelete {
		props, err := sc.command.ossGetObjectStatRetry(bucket, object)
		if err != nil {
			return err
		}
		allheaders = sc.mergeHeader(props, headers, isUpdate, isDelete)
	}

	options, err := sc.getOSSOptions(allheaders)
	if err != nil {
		return err
	}

	return sc.ossSetObjectMetaRetry(bucket, object, options...)
}

func (sc *SetMetaCommand) mergeHeader(props http.Header, headers map[string]string, isUpdate, isDelete bool) map[string]string {
	allheaders := map[string]string{}
	for name := range props {
		if _, err := fetchHeaderOptionMap(name); err == nil || strings.HasPrefix(strings.ToLower(name), strings.ToLower(oss.HTTPHeaderOssMetaPrefix)) {
			allheaders[strings.ToLower(name)] = props.Get(name)
		}
		if name == StatACL {
			allheaders[strings.ToLower(oss.HTTPHeaderOssObjectACL)] = props.Get(name)
		}
	}
	if isUpdate {
		for name, val := range headers {
			allheaders[strings.ToLower(name)] = val
		}
	}
	if isDelete {
		for name := range headers {
			delete(allheaders, strings.ToLower(name))
		}
	}
	return allheaders
}

func (sc *SetMetaCommand) getOSSOptions(headers map[string]string) ([]oss.Option, error) {
	options := []oss.Option{}
	for name, value := range headers {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(oss.HTTPHeaderOssMetaPrefix)) {
			options = append(options, oss.Meta(name[len(oss.HTTPHeaderOssMetaPrefix):], value))
		} else {
			option, err := getOSSOption(name, value)
            if err != nil {
				return nil, err
			}
			options = append(options, option)
		}
	}
	return options, nil
}

func (sc *SetMetaCommand) ossSetObjectMetaRetry(bucket *oss.Bucket, object string, options ...oss.Option) error {
	retryTimes, _ := GetInt(OptionRetryTimes, sc.command.options)
	for i := 1; ; i++ {
		_, err := bucket.CopyObject(object, object, options...)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, object}
		}
	}
}

func (sc *SetMetaCommand) batchSetObjectMeta(bucket *oss.Bucket, cloudURL CloudURL, headers map[string]string, isUpdate, isDelete, force bool, routines int64) error {
	if !force {
		var val string
		fmt.Printf("Do you really mean to recursivlly set meta on objects of %s(y or N)? ", sc.command.args[0])
		if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
			fmt.Println("operation is canceled.")
			return nil
		}
	}

	// producer list objects
	// consumer set meta
	chObjects := make(chan string, ChannelBuf)
	chFinishObjects := make(chan string, ChannelBuf)
	chError := make(chan error, routines+1)
	go sc.command.objectProducer(bucket, cloudURL, chObjects, chError)
	for i := 0; int64(i) < routines; i++ {
		go sc.setObjectMetaConsumer(bucket, headers, isUpdate, isDelete, chObjects, chFinishObjects, chError)
	}

	completed := 0
	num := 0
	for int64(completed) <= routines {
		select {
		case <-chFinishObjects:
			num++
			fmt.Printf("\rsetted object meta on %d objects, when error happens...", num)
		case err := <-chError:
			if err != nil {
				fmt.Printf("\rsetted object meta on %d objects, when error happens.\n", num)
				return err
			}
			completed++
		}
	}
	fmt.Printf("\rSucceed:scanned %d objects, setted object meta on %d objects.\n", num, num)
	return nil
}

func (sc *SetMetaCommand) setObjectMetaConsumer(bucket *oss.Bucket, headers map[string]string, isUpdate, isDelete bool, chObjects <-chan string, chFinishObjects chan<- string, chError chan<- error) {
	for object := range chObjects {
		err := sc.setObjectMeta(bucket, object, headers, isUpdate, isDelete)
		if err != nil {
			chError <- err
			return
		}
		chFinishObjects <- object
	}

	chError <- nil
}
