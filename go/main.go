package main

import (
    // "errors"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/mattn/go-sqlite3"

    "github.com/qor/admin"
    "github.com/qor/l10n"
    "github.com/qor/media"
    "github.com/qor/publish2"
    "github.com/qor/sorting"
    "github.com/qor/validations"

    // "github.com/qor/media"
    "github.com/qor/media/media_library"
    //  "github.com/qor/media/oss"
    // "github.com/qor/auth/providers/facebook"
    // "github.com/qor/auth/providers/github"
    // "github.com/qor/auth/providers/google"
    // "github.com/qor/auth/providers/twitter"
)

// DB Global DB connection
var DB *gorm.DB

type SMTPConfig struct {
    Host     string
    Port     string
    User     string
    Password string
}

var Config = struct {
    HTTPS bool `default:"false" env:"HTTPS"`
    Port  uint `default:"7000" env:"PORT"`
    DB    struct {
        Name     string `env:"DBName" default:"qor_example"`
        Adapter  string `env:"DBAdapter" default:"mysql"`
        Host     string `env:"DBHost" default:"localhost"`
        Port     string `env:"DBPort" default:"3306"`
        User     string `env:"DBUser"`
        Password string `env:"DBPassword"`
    }
    S3 struct {
        AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
        SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
        Region          string `env:"AWS_Region"`
        S3Bucket        string `env:"AWS_Bucket"`
    }
    AmazonPay struct {
        MerchantID   string `env:"AmazonPayMerchantID"`
        AccessKey    string `env:"AmazonPayAccessKey"`
        SecretKey    string `env:"AmazonPaySecretKey"`
        ClientID     string `env:"AmazonPayClientID"`
        ClientSecret string `env:"AmazonPayClientSecret"`
        Sandbox      bool   `env:"AmazonPaySandbox"`
        CurrencyCode string `env:"AmazonPayCurrencyCode" default:"JPY"`
    }
    SMTP SMTPConfig
}{}

/*
func init() {
    var err error

    dbConfig := Config.DB
    if config.Config.DB.Adapter == "mysql" {
        DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
        // DB = DB.Set("gorm:table_options", "CHARSET=utf8")
    } else if config.Config.DB.Adapter == "postgres" {
        DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name))
    } else if config.Config.DB.Adapter == "sqlite" {
        DB, err = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", os.TempDir(), dbConfig.Name))
    } else {
        panic(errors.New("not supported database adapter"))
    }
        DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
    if err == nil {
        if os.Getenv("DEBUG") != "" {
            DB.LogMode(true)
        }

        l10n.RegisterCallbacks(DB)
        sorting.RegisterCallbacks(DB)
        validations.RegisterCallbacks(DB)
        media.RegisterCallbacks(DB)
        publish2.RegisterCallbacks(DB)
    } else {
        panic(err)
    }
}
*/

type Places []Place

type Twint struct {
    gorm.Model
    ConversationID string `json:"conversation_id"`
    // CreatedAt      float64   `json:"created_at"`
    // Date          string    `json:"date"`
    Hashtags      []Hashtag `gorm:"many2many:twint_hashtags;" json:"hashtags"`
    ID            float64   `json:"id"`
    LikesCount    int       `json:"likes_count"`
    Link          string    `json:"link"`
    Location      string    `json:"location"`
    Mentions      string    `json:"mentions"`
    Name          string    `json:"name"`
    Photos        []Photo   `json:"photos"`
    Place         Place     `json:"place"`
    QuoteURL      string    `json:"quote_url"`
    RepliesCount  int       `json:"replies_count"`
    Retweet       Retweet   `json:"retweet"`
    RetweetsCount int       `json:"retweets_count"`
    Time          string    `json:"time"`
    Timezone      string    `json:"timezone"`
    Tweet         string    `json:"tweet"`
    // Urls          []Url     `json:"urls"`
    UserID     float64  `json:"user_id"`
    Username   string   `json:"username"`
    Video      int      `json:"video"`
    CategoryID uint     `l10n:"sync"`
    Category   Category `l10n:"sync"`
}

type Category struct {
    gorm.Model
    l10n.Locale
    sorting.Sorting
    Name string
    Code string

    Categories []Category
    CategoryID uint
}

func (category Category) Validate(db *gorm.DB) {
    if strings.TrimSpace(category.Name) == "" {
        db.AddError(validations.NewError(category, "Name", "Name can not be empty"))
    }
}

func (category Category) DefaultPath() string {
    if len(category.Code) > 0 {
        return fmt.Sprintf("/category/%s", category.Code)
    }
    return "/"
}

type Retweet struct {
    gorm.Model
    ConversationID string `json:"conversation_id"`
}

type Photo struct {
    gorm.Model
    Link string `json:"link"`
}

type Hashtag2 struct {
    gorm.Model
    Label string `json:"label"`
}

type Place struct {
    gorm.Model
    Name string `json:"name"`
}

type Url struct {
    gorm.Model
    Href   string `json:"href"`
    Status int8   `json:"status"`
}

type Profile struct {
    gorm.Model
    BackgroundImage string `json:"background_image"`
    Bio             string `json:"bio"`
    Followers       int    `json:"followers"`
    Following       int    `json:"following"`
    // ID              float64 `json:"id"`
    JoinDate        string `json:"join_date"`
    JoinTime        string `json:"join_time"`
    Likes           int    `json:"likes"`
    Location        string `json:"location"`
    Media           int    `json:"media"`
    Name            string `json:"name"`
    Private         int    `json:"private"`
    ProfileImageURL string `json:"profile_image_url"`
    Tweets          int    `json:"tweets"`
    URLs            []Url  `json:"url"`
    Username        string `json:"username"`
    Verified        int    `json:"verified"`
    // Twints []Twint `json:"tweets"`
}

// Create a GORM-backend model
type User struct {
    gorm.Model
    Name            string
    Username        string
    Bio             string
    Location        string
    URL             string
    JoinDatetime    time.Time
    JoinDate        time.Time
    JoinTime        time.Time
    TweetsCount     int
    FollowingCount  int
    FollowersCount  int
    LikesCount      int
    MediaCount      int
    Private         bool
    Verified        bool
    Avatar          string
    BackgroundImage string
    Session         string
    GeoUser         string
    Follows         []Friend
    Friends         []Friend
    Tweets          []Tweet
}

type Friend struct {
    gorm.Model
    User   string
    Follow string
    Essid  string
}

type Tweet struct {
    gorm.Model
    ConversationId int64
    PublishedAt    int64
    Date           time.Time
    TimeZone       string
    Place          string
    Location       string
    Tweet          string
    Hashtags       []Hashtag `gorm:"many2many:twint_hashtags;" json:"hashtags"`
    // Hashtags        []Hashtag
    UserId    int64
    UserIdStr string
    Username  string
    Name      string
    // ProfileImageUrl string
    ProfileImageUrl media_library.MediaBox
    Day             int
    Hour            int
    Link            string
    Retweet         bool
    Essid           string
    NLikes          int
    NReplies        int
    NRetweets       int
    QuoteURL        string
    Video           bool
    Search          string
    Near            string
    GeoNear         string
    GeoTweet        string
    // Photos          string
    Photos   media_library.MediaBox
    Mentions string
}

type Hashtag struct {
    gorm.Model
    Name    string
    Indices int
    Code    string `l10n:"sync"`
}

func (h Hashtag) Validate(db *gorm.DB) {
    if strings.TrimSpace(h.Name) == "" {
        db.AddError(validations.NewError(h, "Name", "Name can not be empty"))
    }
    // if strings.TrimSpace(h.Code) == "" {
    //    db.AddError(validations.NewError(h, "Code", "Code can not be empty"))
    // }
}

type Image struct {
    ID          int       `json:"id"`          //ID
    SourceUrl   string    `json:"source_url"`  //
    Path        string    `json:"path"`        //
    ReadNum     int       `json:"read_num"`    //
    LikeNum     int       `json:"like_num"`    //
    CommentNum  int       `json:"comment_num"` //
    PublishedAt time.Time `json:"published_at" gorm:"default: null"`
    // CreatedAt   time.Time `json:"created_at"` //
    // UpdatedAt   time.Time `json:"updated_at"` //
}

func main() {
    var err error
    DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", "root", "da33ve79T!", "127.0.0.1", 3306, "twint"))
    if err == nil {
        if os.Getenv("DEBUG") != "" {
            DB.LogMode(true)
        }

        l10n.RegisterCallbacks(DB)
        sorting.RegisterCallbacks(DB)
        validations.RegisterCallbacks(DB)
        media.RegisterCallbacks(DB)
        publish2.RegisterCallbacks(DB)
    } else {
        panic(err)
    }

    // DB, _ := gorm.Open("sqlite3", "tweets.db")
    // DB, _ := gorm.Open("mysql", "tweets.db")
    DB.AutoMigrate(&User{}, &Friend{}, &Hashtag{}, &Url{}, &Tweet{}, &Twint{}, &Image{}, &Profile{}, &Place{}, &Photo{})

    // Initalize
    Admin := admin.New(&admin.AdminConfig{DB: DB})

    // Allow to use Admin to manage User, Product
    Admin.AddResource(&User{})
    Admin.AddResource(&Friend{})
    Admin.AddResource(&Hashtag{})
    Admin.AddResource(&Tweet{})
    Admin.AddResource(&Image{})
    Admin.AddResource(&Profile{})
    Admin.AddResource(&Twint{})
    // Admin.AddResource(&Place{})
    Admin.AddResource(&Url{})
    Admin.AddResource(&Photo{})

    // initalize an HTTP request multiplexer
    mux := http.NewServeMux()

    // Mount admin interface to mux
    Admin.MountTo("/admin", mux)

    fmt.Println("Listening on: 9000")
    http.ListenAndServe(":9000", mux)
}
