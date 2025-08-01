package v1

import (
	"archiveFiles/internal/entity"
	"archiveFiles/internal/httpcontroller/v1/api"
	"time"
)

func convertTask(task *entity.Task) *api.Task {
	now := time.Now()
	id := int(task.ID)
	apiTask := &api.Task{
		Status:  &task.Status,
		Id:      &id,
		Created: &now,
	}

	if len(task.Links) != 0 {
		apiTask.Links = &api.ArchLinks{}
		for _, v := range task.Links {
			archLink := &api.ArchLink{
				Link: &v.URL,
			}
			if v.Err != nil {
				archErr := v.Err.Error()
				archLink.Error = &archErr
			}

			if v.FileName != "" {
				archLink.FileName = &v.FileName
			}

			*apiTask.Links = append(*apiTask.Links, *archLink)
		}
	}

	if task.DownloadLink != "" {
		apiTask.Download = &task.DownloadLink
	}
	if task.Err != nil {
		errMess := task.Err.Error()
		apiTask.Error = &errMess
	}

	return apiTask
}

func convertLinks(links *api.Links) []entity.Link {

	taskLinks := []entity.Link{}
	if links != nil && len(*links) != 0 {

		for _, v := range *links {
			taskLinks = append(taskLinks, entity.Link{URL: v})
		}

	}

	return taskLinks
}
