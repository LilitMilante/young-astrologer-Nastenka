package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"young-astrologer-Nastenka/internal/entity"
)

type Repository interface {
	APODs(ctx context.Context) ([]entity.APOD, error)
	APODByDate(ctx context.Context, date string) (image entity.APOD, err error)
	CreateAPOD(ctx context.Context, apod entity.APOD) error
}

type Service struct {
	nasaAPI string
	repo    Repository
	client  *http.Client
}

func NewService(nasaAPI string, repo Repository) *Service {
	return &Service{
		nasaAPI: nasaAPI,
		repo:    repo,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (s *Service) APODs(ctx context.Context) ([]entity.APOD, error) {
	return s.repo.APODs(ctx)
}

func (s *Service) APODByDate(ctx context.Context, date string) (image entity.APOD, err error) {
	return s.repo.APODByDate(ctx, date)
}

func (s *Service) UpdateAPODJob(ctx context.Context) {
	updateTime := time.Now()

	timer := time.NewTimer(0) // 0 to run on start once
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			err := s.updateAPOD(ctx)
			if err != nil {
				log.Printf("update APOD: %s", err)
			}

			// Setting the timer for the next day at 00:01
			updateTime = updateTime.Add(24 * time.Hour)
			timer.Reset(updateTime.Sub(time.Now()))
		}
	}
}

func (s *Service) updateAPOD(ctx context.Context) error {
	today := time.Now().UTC().Format("2006-01-02")

	apod, err := s.repo.APODByDate(ctx, today)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("get apod by %s: %w", today, err)
	}

	apod, err = s.downloadAPOD(ctx)
	if err != nil {
		return fmt.Errorf("dowload APOD: %w", err)
	}

	apod.ImageB64, err = s.downloadImage(ctx, apod.Url)
	if err != nil {
		return fmt.Errorf("download image by url: %w", err)
	}

	return s.repo.CreateAPOD(ctx, apod)
}

func (s *Service) downloadAPOD(ctx context.Context) (apod entity.APOD, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.nasaAPI, nil)
	if err != nil {
		return apod, fmt.Errorf("create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return apod, fmt.Errorf("send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apod, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&apod)
	if err != nil {
		return apod, fmt.Errorf("decode request body: %w", err)
	}

	return apod, nil
}

// downloadImage downloads an image from url and returns it in base64 format.
func (s *Service) downloadImage(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
