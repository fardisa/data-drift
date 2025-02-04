package history

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/reducers"
	"github.com/data-drift/data-drift/urlgen"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v56/github"
	"github.com/shopspring/decimal"
)

func ProcessHistory(client *github.Client, redisClient *redis.Client, repoOwner string, repoName string, metric common.MetricConfig, installationId int) (common.MetricStorageKey, error) {
	kpiRepository := common.NewKpiRepository(redisClient)
	reportBaseUrl := urlgen.BuildReportDiffBaseUrl(repoOwner, repoName)
	fmt.Println(reportBaseUrl)
	ctx := context.Background()

	csvFilePath := metric.Filepath
	dateColumnName := metric.DateColumnName
	KPIColumnName := metric.KPIColumnName
	metricName := metric.MetricName

	fmt.Println(metric)
	// Set the start and end dates to display the history for.
	endDate := time.Now()

	if dateColumnName == "" {
		dateColumnName = "date"
	}

	// Get the commit history for the repository.
	// Get the commit history for the file.
	commits, _, err := client.Repositories.ListCommits(context.Background(), repoOwner, repoName, &github.CommitsListOptions{
		Path:        csvFilePath,
		SHA:         "",
		Until:       endDate,
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return "", fmt.Errorf("error getting commit history: %v", err.Error())
	}

	// Print the number of commits.
	fmt.Printf("Number of commits: %d\n", len(commits))

	// Group the lines of the CSV file by reporting date.
	lineCountAndKPIByDateByVersion := make(common.Metrics)
	for index, commit := range commits {
		var commitMessages []common.CommitComments
		fmt.Printf("\r Commit %d/%d", index, len(commits))

		commitSha := common.CommitSha(*commit.SHA)

		commitDate := commit.Commit.Committer.Date
		commitComments := GetCommitComments(client, ctx, repoOwner, repoName, *commit.SHA)

		for _, comment := range commitComments {
			commitMessages = append(commitMessages, common.CommitComments{CommentBody: *comment.Body, CommentAuthor: *comment.User.Login})
		}
		commitTimestamp := commitDate.Unix()
		records, err := getFileContentsForCommit(client, repoOwner, repoName, csvFilePath, *commit.SHA)
		if err != nil {
			log.Printf("Error getting file contents for commit %s: %v", *commit.SHA, err.Error())
			continue
		}

		if err != nil {
			log.Printf("Error parsing CSV file for commit %s: %v", *commit.SHA, err.Error())
			continue
		}
		var dateColumn int
		var kpiColumn int
		var defaultDateColumn int
		var dimensionColumns []struct {
			dimensionName   string
			dimensionColumn int
		}

		for i, columnName := range records[0] {
			if columnName == dateColumnName {
				dateColumn = i
			}
			if columnName == "date" {
				defaultDateColumn = i
			}
			if columnName == KPIColumnName {
				kpiColumn = i
			}
			for _, metricDimension := range metric.Dimensions {
				if columnName == metricDimension {
					dimensionColumns = append(dimensionColumns, struct {
						dimensionName   string
						dimensionColumn int
					}{
						dimensionName:   metricDimension,
						dimensionColumn: i,
					})
				}
			}
		}

		for _, record := range records[1:] { // Skip the header row.
			for _, timegrain := range GetDefaultTimeGrains(metric.TimeGrains) {
				var periodKey common.PeriodKey
				dateValue := record[dateColumn]
				if len(dateValue) > 10 {
					dateValue = dateValue[:10]
				}
				periodTime, parsingError := time.Parse("2006-01-02", dateValue)

				if parsingError != nil {
					dateValue := record[defaultDateColumn]
					if len(dateValue) > 10 {
						dateValue = dateValue[:10]
					}
					periodTime, parsingError = time.Parse("2006-01-02", dateValue)
					if parsingError != nil {
						fmt.Println("Error with default date:", parsingError.Error())
						continue
					}
				}

				switch timegrain {
				case common.Day:
					periodKey = common.PeriodKey(periodTime.Format("2006-01-02"))
				case common.Week:
					_, week := periodTime.ISOWeek()
					periodKey = common.PeriodKey(fmt.Sprintf("%d-W%02d", periodTime.Year(), week))
				case common.Month:
					periodKey = common.PeriodKey(periodTime.Format("2006-01"))
				case common.Quarter:
					periodKey = common.PeriodKey(fmt.Sprintf("%d-Q%d", periodTime.Year(), (periodTime.Month()-1)/3+1))
				case common.Year:
					periodKey = common.PeriodKey(periodTime.Format("2006"))
				default:
					fmt.Printf("Invalid time grain: %s", timegrain)
				}

				periodAndDimensionKey := common.PeriodAndDimensionKey(string(periodKey))
				dimension := common.Dimension("none")
				dimensionValue := common.DimensionValue(common.NoDimensionValue)

				updateMetric(lineCountAndKPIByDateByVersion, periodAndDimensionKey, timegrain, periodKey, dimension, dimensionValue, record, kpiColumn, commitSha, commitTimestamp, commit, commitMessages, reportBaseUrl)

				for _, metricDimension := range dimensionColumns {
					dimension = common.Dimension(metricDimension.dimensionName)
					dimensionValue = common.DimensionValue(record[metricDimension.dimensionColumn])
					periodAndDimensionKey = common.PeriodAndDimensionKey(string(periodKey) + " " + string(dimensionValue))
					updateMetric(lineCountAndKPIByDateByVersion, periodAndDimensionKey, timegrain, periodKey, dimension, dimensionValue, record, kpiColumn, commitSha, commitTimestamp, commit, commitMessages, reportBaseUrl)

				}
			}
		}

	}

	// Print the line count for each reporting date.
	for dateStr, lineCounts := range lineCountAndKPIByDateByVersion {

		var countsStr string
		var kpiStr string
		for _, count := range lineCounts.History {
			countsStr += fmt.Sprintf("%d ", count.Lines)
		}
		fmt.Printf("Line Count %s: %s\n", dateStr, countsStr)
		fmt.Printf("       KPI %s: %s\n", dateStr, kpiStr)
	}

	if _, err := os.Stat("dist"); os.IsNotExist(err) {
		if err := os.Mkdir("dist", 0755); err != nil {
			fmt.Printf("Error creating directory: %v", err.Error())
		}
	}

	// Generate a timestamp to include in the JSON file name.
	// Open a file to write the line counts by date by version in JSON format.
	// Write the line counts and KPI values to the JSON file.
	metricStoredFilePath := kpiRepository.WriteMetricKPI(repoOwner, repoName, metricName, lineCountAndKPIByDateByVersion)
	return metricStoredFilePath, nil
}

func updateMetric(lineCountAndKPIByDateByVersion common.Metrics, periodAndDimensionKey common.PeriodAndDimensionKey, timegrain common.TimeGrain, periodKey common.PeriodKey, dimension common.Dimension, dimensionValue common.DimensionValue, record []string, kpiColumn int, commitSha common.CommitSha, commitTimestamp int64, commit *github.RepositoryCommit, commitMessages []common.CommitComments, reportBaseUrl string) {
	if lineCountAndKPIByDateByVersion[periodAndDimensionKey].History == nil {
		lineCountAndKPIByDateByVersion[periodAndDimensionKey] = common.Metric{
			TimeGrain:      timegrain,
			Period:         periodKey,
			Dimension:      dimension,
			DimensionValue: dimensionValue,
			History:        make(map[common.CommitSha]common.CommitData),
		}
	}

	kpiStr := record[kpiColumn]
	kpi, _ := decimal.NewFromString(kpiStr)

	newLineCount := lineCountAndKPIByDateByVersion[periodAndDimensionKey].History[commitSha].Lines + 1
	newKPI := kpi.Add(lineCountAndKPIByDateByVersion[periodAndDimensionKey].History[commitSha].KPI)

	firstDateOfPeriod, _ := reducers.LegacyGetFirstComputationDateOfPeriod(periodKey)
	isAfterPeriod := time.Unix(commitTimestamp, 0).After(firstDateOfPeriod)
	urlQueryStringForCommitUrl, _ := reducers.GetQueryStringFiltersForPeriod(periodKey, dimension, dimensionValue)

	lineCountAndKPIByDateByVersion[periodAndDimensionKey].History[commitSha] = common.CommitData{
		Lines:           newLineCount,
		KPI:             newKPI,
		CommitTimestamp: commitTimestamp,
		CommitDate:      time.Unix(commitTimestamp, 0).Format("2006-01-02"),
		IsAfterPeriod:   isAfterPeriod,
		CommitUrl:       urlgen.BuildReportDiffUrl(reportBaseUrl, string(commitSha), urlQueryStringForCommitUrl),
		CommitComments:  commitMessages,
	}
}

func getFileContentsForCommit(client *github.Client, owner, name, path, sha string) ([][]string, error) {
	opts := &github.RepositoryContentGetOptions{Ref: sha}
	fileContents, _, ghresp, err := client.Repositories.GetContents(context.Background(), owner, name, path, opts)
	if err != nil {
		return nil, err
	}
	downloadableUrl := fileContents.GetDownloadURL()

	resp, err := http.Get(downloadableUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	defer ghresp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()

	return records, nil
}

func GetDefaultTimeGrains(timeGrains []common.TimeGrain) []common.TimeGrain {
	if len(timeGrains) == 0 {
		return []common.TimeGrain{common.Month}
	}
	return timeGrains
}
