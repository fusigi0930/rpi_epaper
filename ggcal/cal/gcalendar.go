package cal

import (
	"os"
	"fmt"
	"io/ioutil"
	//"net/http"
	"context"
	"time"
	"runtime"

	"ggcal/log"

	//"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type WeekDay int

const (
	Sun WeekDay = iota
	Mon	
	Tue
	Wed
	Thu
	Fri
	Sat
	Unknow
)

type CalConfig struct {
	CredFile string `yaml:"cred_file"`
	HolidayCal string `yaml:"holiday_cal"`
	Calendar []string `yaml:"calendar"`
}

var gCal CalConfig
var gCalService *calendar.Service

func LoadDef(file string) error {
	var config_folder string
	var fullname string
	switch runtime.GOOS {
	case "windows":
		config_folder = fmt.Sprintf("%s/ggcal", os.Getenv("ProgramData"))
		fullname = fmt.Sprintf("%s/%s", config_folder, file)
	case "linux":
		config_folder = "/etc/ggcal"
		fullname = fmt.Sprintf("%s/%s", config_folder, file)
	}
	log.LogService().Infof("get config path: %s\n", config_folder)

	f, err := os.ReadFile(fullname)
	if err != nil {
		log.LogService().Errorf("load calendar config failed: %v\n", err)
		return err
	}

	if err := yaml.Unmarshal(f, &gCal); err != nil {
		log.LogService().Errorf("parsing config failed: %v\n", err)
		return err	
	}

	cred_file := fmt.Sprintf("%s/%s", config_folder, gCal.CredFile)
	gCalService, err = getGoogleCalendarService(cred_file, "")
	if err != nil {
		log.LogService().Errorf("get google calendar service failed: %v\n", err)
		return err			
	}

	return nil
}

func CalService() *calendar.Service {
	if gCalService == nil {
		log.LogService().Errorf("calendar service still null\n")
	}
	return gCalService
}

func getGoogleCalendarService(keyFile string, calId string) (*calendar.Service, error) {
	ctx := context.Background()

	jsonKey, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.LogService().Errorf("failed to read service key file: %v\n", err)
		return nil, err
	}

	jwtConfig, err := google.JWTConfigFromJSON(jsonKey, calendar.CalendarReadonlyScope)
	if err != nil {
		log.LogService().Errorf("failed to parse key file to jwt config: %v\n", err)
		return nil, err
	}

	if len(calId) != 0 {
		jwtConfig.Subject = calId
	}

	client := jwtConfig.Client(ctx)
	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.LogService().Errorf("get calendar client and service failed: %v\n", err)
		return nil, err
	}

	return service, nil
}

//func ListEvents(service *calendar.Service, calId string, maxResult int64) error {
//	t := time.Now().Format(time.RFC3339)
//	events, err := service.Events.List(calId).ShowDeleted(false).SingleEvents(true).
//				TimeMin(t).MaxResults(maxResults).OrderBy("startTime").Do()
//
//	if err != nil {
//		log.LogService().Errorf("get events object filed: %v\n", err)
//		return err
//	}
//
//	return nil
//}

func GetWeekDay(date string) WeekDay {
	if len(date) > 10 {
		date = date[:10]
	}

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return Unknow
	}

	return WeekDay(t.Weekday())
}

func GetMonthDayCount(date string) int {
	if len(date) > 10 {
		date = date[:10]
	}

	this_mon := fmt.Sprintf("%s-01", date[:len(date)-3])
	t, err := time.Parse("2006-01-02", this_mon)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return 0
	}

	next_mon := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	day_count := next_mon.Add(-24 * time.Hour).Day()

	return day_count
}

func IsHoliday(service *calendar.Service, date string) (bool, string) {
	holidayId := gCal.HolidayCal

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return false, ""
	}

	next_t := t.Add(24 * time.Hour)

	events, err := service.Events.List(holidayId).
				TimeMin(t.Format(time.RFC3339)).
				TimeMax(next_t.Format(time.RFC3339)).
				SingleEvents(true).Do()

	if err != nil {
		log.LogService().Errorf("calendar service error: %v\n", err)
		return false, ""
	}

	if len(events.Items) > 0 {
		return true, events.Items[0].Summary
	}

	return false, ""
}

func GetEventsFromDate(service *calendar.Service, date string) ([]*calendar.Event, error) {
	calId := gCal.Calendar[0]
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return []*calendar.Event{}, err
	}

	next_t := t.Add(24 * time.Hour)

	events, err := service.Events.List(calId).
				ShowDeleted(false).
				TimeMin(t.Format(time.RFC3339)).
				TimeMax(next_t.Format(time.RFC3339)).
				OrderBy("startTime").
				SingleEvents(true).Do()	

	if err != nil {
		log.LogService().Errorf("calendar service error: %v\n", err)
		return []*calendar.Event{}, err
	}

	if len(events.Items) == 0 {
		return []*calendar.Event{}, nil
	}

	return events.Items, nil
}

func GetHolidays(service *calendar.Service, date string) map[string]string {
	holidayId := gCal.HolidayCal

	dayCount := GetMonthDayCount(date)
	startTime := fmt.Sprintf("%s-01 00:00:00", date[:len(date)-3])
	endTime := fmt.Sprintf("%s-%02d 23:59:59", date[:len(date)-3], dayCount)

	st, err := time.Parse("2006-01-02 15:04:05",  startTime)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return nil
	}

	et, err := time.Parse("2006-01-02 15:04:05",  endTime)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return nil
	}

	events, err := service.Events.List(holidayId).
				TimeMin(st.Format(time.RFC3339)).
				TimeMax(et.Format(time.RFC3339)).
				SingleEvents(true).Do()

	if err != nil {
		log.LogService().Errorf("calendar service error: %v\n", err)
		return nil
	}

	holidays := make(map[string]string, 0)
	for _, e := range events.Items {
		holidays[e.Start.Date] = e.Summary
	}

	return holidays
}

func GetEvents(service *calendar.Service, date string) map[string]string {
	calId := gCal.Calendar[0]

	dayCount := GetMonthDayCount(date)
	startTime := fmt.Sprintf("%s-01 00:00:00", date[:len(date)-3])
	endTime := fmt.Sprintf("%s-%02d 23:59:59", date[:len(date)-3], dayCount)

	st, err := time.Parse("2006-01-02 15:04:05",  startTime)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return nil
	}

	et, err := time.Parse("2006-01-02 15:04:05",  endTime)
	if err != nil {
		log.LogService().Errorf("wrong date string: %v\n", err)
		return nil
	}

	events, err := service.Events.List(calId).
				ShowDeleted(false).
				TimeMin(st.Format(time.RFC3339)).
				TimeMax(et.Format(time.RFC3339)).
				OrderBy("startTime").
				SingleEvents(true).Do()	

	if err != nil {
		log.LogService().Errorf("calendar service error: %v\n", err)
		return nil
	}

	es := make(map[string]string, 0)
	for _, e := range events.Items {
	
		t, _ := time.Parse("2006-01-02T15:04:05-07:00", e.Start.DateTime)
		date := t.Format("2006-01-02")
		el, ok := es[date]

		if !ok {
			es[date] = fmt.Sprintf("%s %s", t.Format("15:04"), e.Summary)
		} else {
			es[date] = fmt.Sprintf("%s\n%s %s", el, t.Format("15:04"), e.Summary)
		}
	}

	return es
}

func GetDaysEvents(service *calendar.Service, date string, days int) map[string]string {
	prefix := []string{"* ", "+ ", "- ", "~ "}
	es := make(map[string]string, 0)
	for i, calId := range gCal.Calendar {
		if i > len(prefix) {
			break
		}
		st, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.LogService().Errorf("wrong date string: %v\n", err)
			return nil
		}

		et := st.Add(time.Duration(days) * 24 * time.Hour)

		events, err := service.Events.List(calId).
					ShowDeleted(false).
					TimeMin(st.Format(time.RFC3339)).
					TimeMax(et.Format(time.RFC3339)).
					OrderBy("startTime").
					SingleEvents(true).Do()	

		if err != nil {
			log.LogService().Errorf("calendar service error: %v\n", err)
			return nil
		}

		for _, e := range events.Items {	
			t, _ := time.Parse("2006-01-02T15:04:05-07:00", e.Start.DateTime)
			date := t.Format("2006-01-02")
			log.LogService().Debugf("**** %s-%s\n", date, e.Summary)
			el, ok := es[date]
			if !ok {
				es[date] = fmt.Sprintf("%s%s %s", prefix[i], t.Format("15:04"), e.Summary)
			} else {
				es[date] = fmt.Sprintf("%s\n%s%s %s", el, prefix[i], t.Format("15:04"), e.Summary)
			}
		}
	}

	return es
}