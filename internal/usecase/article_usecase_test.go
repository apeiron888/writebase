package usecase_test

import (
	"context"
	"testing"
	"write_base/internal/domain"
	"write_base/internal/mocks"
	"write_base/internal/usecase"

	"github.com/stretchr/testify/mock"
)

func TestArticleUsecase_CreateArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		input     *domain.Article
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:   "success",
			userID: "user1",
			input: &domain.Article{
				Title:    "Test",
				AuthorID: "user1",
				Tags:     []string{"go"},
				ContentBlocks: []domain.ContentBlock{
					{Type: "paragraph", Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "Hello"}}},
				},
				Excerpt: "Short",
			},
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				repo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.Article")).
					Return(&domain.Article{ID: "article1"}, nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.CreateArticle(context.Background(), tc.userID, tc.input)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_UpdateArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		input     *domain.Article
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:   "success",
			userID: "user1",
			input: &domain.Article{
				ID:       "article1",
				Title:    "Updated Test",
				Slug:     "updated-test",
				AuthorID: "user1",
				Tags:     []string{"go", "test"},
				Language: "en",
				Status:   domain.StatusDraft,
				Excerpt:  "Updated Short",
				ContentBlocks: []domain.ContentBlock{
					{Type: "paragraph", Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "Updated Hello"}}},
				},
			},
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)

				// Simulate old article for comparison
				oldArticle := &domain.Article{
					ID:       "article1",
					Title:    "Old Title",
					Slug:     "old-title",
					AuthorID: "user1",
					Tags:     []string{"go"},
					Language: "en",
					Status:   domain.StatusDraft,
					Excerpt:  "Old excerpt",
					ContentBlocks: []domain.ContentBlock{
						{Type: "paragraph", Content: domain.BlockContent{Paragraph: &domain.ParagraphContent{Text: "Old content"}}},
					},
				}

				repo.On("GetByID", mock.Anything, "article1").Return(oldArticle, nil)
				policy.On("UserOwnsArticle", "user1", *oldArticle).Return(true)

				repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Article")).
					Return(&domain.Article{ID: "article1"}, nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.UpdateArticle(context.Background(), tc.userID, tc.input)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_DeleteArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "user1",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				policy.On("UserOwnsArticle", "user1", domain.Article{ID: "article1"}).Return(true)
				repo.On("GetByID", mock.Anything, "article1").Return(&domain.Article{ID: "article1"}, nil)
				repo.On("Delete", mock.Anything, "article1").Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			err := uc.DeleteArticle(context.Background(), tc.userID, tc.articleID)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_RestoreArticle(t *testing.T) {
    type testCase struct {
        name      string
        userID    string
        articleID string
        setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
        expectErr error
    }

    tests := []testCase{
        {
            name:      "success",
            userID:    "user1",
            articleID: "article1",
            setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
                policy.On("UserExists", "user1").Return(true)
                policy.On("UserOwnsArticle", "user1", domain.Article{ID: "article1"}).Return(true)
                repo.On("GetByID", mock.Anything, "article1").Return(&domain.Article{ID: "article1", Status: domain.StatusDeleted}, nil)
                repo.On("Restore", mock.Anything, "article1").Return(&domain.Article{ID: "article1", Status: domain.StatusDraft}, nil)
            },
            expectErr: nil,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            repo := new(mocks.MockIArticleRepository)
            policy := new(mocks.MockIArticlePolicy)
            if tc.setup != nil {
                tc.setup(repo, policy)
            }
            uc := usecase.NewArticleUsecase(repo, policy)

            _, err := uc.RestoreArticle(context.Background(), tc.userID, tc.articleID)
            if err != tc.expectErr {
                t.Errorf("expected error %v, got %v", tc.expectErr, err)
            }

            repo.AssertExpectations(t)
            policy.AssertExpectations(t)
        })
    }
}

//============================================================================//
//                         Statistics tests                                 //
//============================================================================//

func TestArticleUsecase_GetArticleStats(t *testing.T) {
    type testCase struct {
        name      string
        articleID string
        setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
        expectErr error
    }

    tests := []testCase{
        {
            name:      "success",
            articleID: "article1",
            setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
                repo.On("GetByID", mock.Anything, "article1").Return(&domain.Article{ID: "article1", 
                Stats: domain.ArticleStats{ViewCount: 100, ClapCount: 50}}, nil)
            },
            expectErr: nil,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            repo := new(mocks.MockIArticleRepository)
            policy := new(mocks.MockIArticlePolicy)
            if tc.setup != nil {
                tc.setup(repo, policy)
            }
            uc := usecase.NewArticleUsecase(repo, policy)

            stats, err := uc.GetArticleStats(context.Background(), tc.articleID)
            if err != tc.expectErr {
                t.Errorf("expected error %v, got %v", tc.expectErr, err)
            }
            if stats.ViewCount != 100 || stats.ClapCount != 50 {
                t.Errorf("expected stats ViewCount=100, ClapCount=50; got ViewCount=%d, ClapCount=%d", stats.ViewCount, stats.ClapCount)
            }

            repo.AssertExpectations(t)
            policy.AssertExpectations(t)
        })
    }
}

//============================================================================//
//                         Article State Management tests                       //
//============================================================================//
func TestArticleUsecase_PublishArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "user1",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				policy.On("UserOwnsArticle", "user1", domain.Article{ID: "article1"}).Return(true)
				article := &domain.Article{ID: "article1", Status: domain.StatusDraft}
				repo.On("GetByID", mock.Anything, "article1").Return(article, nil)
				repo.On("Publish", mock.Anything, "article1", mock.AnythingOfType("time.Time")).Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.PublishArticle(context.Background(), tc.userID, tc.articleID)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_UnpublishArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "user1",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				policy.On("UserOwnsArticle", "user1", domain.Article{ID: "article1"}).Return(true)
				article := &domain.Article{ID: "article1", Status: domain.StatusPublished}
				repo.On("GetByID", mock.Anything, "article1").Return(article, nil)
				repo.On("Unpublish", mock.Anything, "article1").Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.UnpublishArticle(context.Background(), tc.userID, tc.articleID, false)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_ArchiveArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "user1",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				policy.On("UserOwnsArticle", "user1", domain.Article{ID: "article1"}).Return(true)
				article := &domain.Article{ID: "article1", Status: domain.StatusPublished}
				repo.On("GetByID", mock.Anything, "article1").Return(article, nil)
				repo.On("Archive", mock.Anything, "article1", mock.AnythingOfType("time.Time")).Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.ArchiveArticle(context.Background(), tc.userID, tc.articleID)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_UnarchiveArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "user1",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				policy.On("UserOwnsArticle", "user1", domain.Article{ID: "article1"}).Return(true)
				article := &domain.Article{ID: "article1", Status: domain.StatusArchived}
				repo.On("GetByID", mock.Anything, "article1").Return(article, nil)
				repo.On("Unarchive", mock.Anything, "article1").Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.UnarchiveArticle(context.Background(), tc.userID, tc.articleID)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

//===========================================================================//
//                         Article Retrieval Tests                                 //
//===========================================================================//
func TestArticleUsecase_GetArticleByID(t *testing.T) {
	type testCase struct {
		name      string
		viewerID  string
		articleID string
		userRole  string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "viewer is author",
			viewerID:  "viewer1",
			articleID: "article1",
			userRole:  "user",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "viewer1").Return(true)
				// article is authored by viewer1
                policy.On("IsAdmin", "viewer1", "user").Return(false)
				article := &domain.Article{
					ID:       "article1",
					AuthorID: "viewer1",
					Status:   domain.StatusDraft,
				}
				repo.On("GetByID", mock.Anything, "article1").Return(article, nil)
			},

			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.GetArticleByID(context.Background(), tc.viewerID, tc.articleID, tc.userRole)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_ViewArticleBySlug(t *testing.T) {
    type testCase struct {
        name      string
        slug      string
        clientIP  string
        setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
        expectErr error
    }

    tests := []testCase{
        {
            name:     "success",
            slug:     "test-article",
            clientIP: "192.168.1.1",   
            setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
                article := &domain.Article{
                    ID:    "article1",
                    Slug:  "test-article",
                    Status: domain.StatusPublished,
                }
                repo.On("GetBySlug", mock.Anything, "test-article").Return(article, nil)
            },
            expectErr: nil,
        },
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            repo := new(mocks.MockIArticleRepository)
            policy := new(mocks.MockIArticlePolicy)
            if tc.setup != nil {
                tc.setup(repo, policy)
            }
            uc := usecase.NewArticleUsecase(repo, policy)

            _, err := uc.ViewArticleBySlug(context.Background(), tc.slug, tc.clientIP)
            if err != tc.expectErr {
                t.Errorf("expected error %v, got %v", tc.expectErr, err)
            }

            repo.AssertExpectations(t)
            policy.AssertExpectations(t)
        })
    } 
}   


//===========================================================================////
//                         Article Listing Tests                                 //
//===========================================================================//
func TestArticleUsecase_ListUserArticles(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		authorID  string
		pag       domain.Pagination
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:     "success",
			userID:   "user1",
			authorID: "author1",
			pag:      domain.Pagination{Page: 1, PageSize: 10, SortField: "created_at", SortOrder: "desc"},
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				repo.On("ListByAuthor", mock.Anything, "author1", mock.AnythingOfType("domain.Pagination")).
					Return([]domain.Article{{ID: "article1"}}, 1, nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			articles, count, err := uc.ListUserArticles(context.Background(), tc.userID, tc.authorID, tc.pag)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
			if count != 1 || len(articles) != 1 || articles[0].ID != "article1" {
				t.Errorf("expected 1 article with ID 'article1', got %d articles with ID '%s'", len(articles), articles[0].ID)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_ListTrendingArticles(t *testing.T) {
	type testCase struct {
		name       string
		pag        domain.Pagination
		windowDays int
		setup      func(*mocks.MockIArticleRepository)
		expectErr  error
		expectDays int
	}

	tests := []testCase{
		{
			name:       "success with custom window",
			pag:        domain.Pagination{Page: 1, PageSize: 5},
			windowDays: 14,
			expectDays: 14,
			setup: func(repo *mocks.MockIArticleRepository) {
				repo.On("ListTrending", mock.Anything, mock.AnythingOfType("domain.Pagination"), 14).
					Return([]domain.Article{{ID: "a1"}}, 1, nil)
			},
			expectErr: nil,
		},
		{
			name:       "defaults to 7-day window",
			pag:        domain.Pagination{Page: 1, PageSize: 5},
			windowDays: 0,
			expectDays: 7,
			setup: func(repo *mocks.MockIArticleRepository) {
				repo.On("ListTrending", mock.Anything, mock.AnythingOfType("domain.Pagination"), 7).
					Return([]domain.Article{{ID: "a1"}}, 1, nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy) // not used here
			if tc.setup != nil {
				tc.setup(repo)
			}

			uc := usecase.NewArticleUsecase(repo, policy)
			articles, count, err := uc.ListTrendingArticles(context.Background(), tc.pag, tc.windowDays)

			if err != tc.expectErr {
				t.Errorf("expected err %v, got %v", tc.expectErr, err)
			}
			if count != 1 || len(articles) != 1 || articles[0].ID != "a1" {
				t.Errorf("unexpected result: got count=%d, articles=%v", count, articles)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_ListArticlesByTag(t *testing.T) {
    type testCase struct {
        name      string
        userID    string
        tag       string
        pag       domain.Pagination
        setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
        expectErr error
    }

    tests := []testCase{
        {
            name:   "success",
            userID: "user1",
            tag:    "golang",
            pag:    domain.Pagination{Page: 1, PageSize: 10},
            setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
                policy.On("UserExists", "user1").Return(true)
                repo.On("ListByTag", mock.Anything, "golang", mock.AnythingOfType("domain.Pagination")).
                    Return([]domain.Article{{ID: "article1"}}, 1, nil)
            },
            expectErr: nil,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            repo := new(mocks.MockIArticleRepository)
            policy := new(mocks.MockIArticlePolicy)
            if tc.setup != nil {
                tc.setup(repo, policy)
            }
            uc := usecase.NewArticleUsecase(repo, policy)

            articles, count, err := uc.ListArticlesByTag(context.Background(), tc.userID, tc.tag, tc.pag)
            if err != tc.expectErr {
                t.Errorf("expected error %v, got %v", tc.expectErr, err)
            }
            if count != 1 || len(articles) != 1 || articles[0].ID != "article1" {
                t.Errorf("expected 1 article with ID 'article1', got %d articles with ID '%s'", len(articles), articles[0].ID)
            }

            repo.AssertExpectations(t)
            policy.AssertExpectations(t)
        })
    }
}

func TestArticleUsecase_SearchArticles(t *testing.T) {
    type testCase struct {
        name      string
        userID    string
        query     string
        pag       domain.Pagination
        setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
        expectErr error
    }

    tests := []testCase{
        {
            name:   "success",
            userID: "user1",
            query:  "golang",
            pag:    domain.Pagination{Page: 1, PageSize: 10},
            setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
                policy.On("UserExists", "user1").Return(true)
                repo.On("Search", mock.Anything, "golang", mock.AnythingOfType("domain.Pagination")).
                    Return([]domain.Article{{ID: "article1"}}, 1, nil)
            },
            expectErr: nil,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            repo := new(mocks.MockIArticleRepository)
            policy := new(mocks.MockIArticlePolicy)
            if tc.setup != nil {
                tc.setup(repo, policy)
            }
            uc := usecase.NewArticleUsecase(repo, policy)

            articles, count, err := uc.SearchArticles(context.Background(), tc.userID, tc.query, tc.pag)
            if err != tc.expectErr {
                t.Errorf("expected error %v, got %v", tc.expectErr, err)
            }
            if count != 1 || len(articles) != 1 || articles[0].ID != "article1" {
                t.Errorf("expected 1 article with ID 'article1', got %d articles with ID '%s'", len(articles), articles[0].ID)
            }

            repo.AssertExpectations(t)
            policy.AssertExpectations(t)
        })
    }
}

func TestArticleUsecase_FilterArticles(t *testing.T) {
	type testcase struct {
		name 	string
		userID 	string
		filter domain.ArticleFilter
		pag 	domain.Pagination
		setup 	func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}
	tests := []testcase{
		{
			name:      "success",
			userID:    "user1",
			filter:    domain.ArticleFilter{Tags: []string{"golang"}},
			pag:      domain.Pagination{Page: 1, PageSize: 10},
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				repo.On("Filter", mock.Anything, domain.ArticleFilter{Tags: []string{"golang"}}, domain.Pagination{Page: 1, PageSize: 10}).Return([]domain.Article{{ID: "article1"}}, 1, nil)

			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, _, err := uc.FilterArticles(context.Background(), tc.userID, tc.filter, tc.pag)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}


//============================================================================//
//                         Engagement Tests                             //
//============================================================================//
func TestArticleUsecase_ClapArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "user1",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				repo.On("GetByID", mock.Anything, "article1").Return(&domain.Article{ID: "article1"}, nil)
				repo.On("IncrementClap", mock.Anything, "article1").Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.ClapArticle(context.Background(), tc.articleID, tc.userID)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

//============================================================================//
//                         Trash Management Tests                             //
//============================================================================//

func TestArticleUsecase_EmptyTrash(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:   "success",
			userID: "user1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "user1").Return(true)
				repo.On("EmptyTrash", mock.Anything, "user1").Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			err := uc.EmptyTrash(context.Background(), tc.userID)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}

			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}
// ===========================================================================//
//	                     Admin Operations Testing                              //
// ===========================================================================//

func TestArticleUsecase_AdminListAllArticles(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		userRole  string
		pag       domain.Pagination
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "admin",
			userRole:  "admin",
			pag:       domain.Pagination{Page: 1, PageSize: 10},
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("UserExists", "admin").Return(true).Once()
				policy.On("IsAdmin", "admin", "admin").Return(true).Once()
				repo.On("AdminListAllArticles", mock.Anything, mock.AnythingOfType("domain.Pagination")).
					Return([]domain.Article{{ID: "article1"}}, 1, nil).Once()
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, _, err := uc.AdminListAllArticles(context.Background(), tc.userID, tc.userRole, tc.pag)
			if err != tc.expectErr {
				t.Fatalf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_AdminHardDeleteArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		userRole  string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "admin",
			userRole:  "admin",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("IsAdmin", "admin", "admin").Return(true).Once()
				repo.On("GetByID", mock.Anything, "article1").Return(&domain.Article{ID: "article1"}, nil).Once()
				repo.On("HardDelete", mock.Anything, "article1").Return(nil).Once()
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			err := uc.AdminHardDeleteArticle(context.Background(), tc.userID, tc.userRole, tc.articleID)
			if err != tc.expectErr {
				t.Fatalf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}

func TestArticleUsecase_AdminUnpublishArticle(t *testing.T) {
	type testCase struct {
		name      string
		userID    string
		userRole  string
		articleID string
		setup     func(*mocks.MockIArticleRepository, *mocks.MockIArticlePolicy)
		expectErr error
	}

	tests := []testCase{
		{
			name:      "success",
			userID:    "admin",
			userRole:  "admin",
			articleID: "article1",
			setup: func(repo *mocks.MockIArticleRepository, policy *mocks.MockIArticlePolicy) {
				policy.On("IsAdmin", "admin", "admin").Return(true).Once()
				repo.On("GetByID", mock.Anything, "article1").
					Return(&domain.Article{ID: "article1", Status: domain.StatusPublished}, nil).Once()
				repo.On("Unpublish", mock.Anything, "article1").Return(nil).Once()
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockIArticleRepository)
			policy := new(mocks.MockIArticlePolicy)
			if tc.setup != nil {
				tc.setup(repo, policy)
			}
			uc := usecase.NewArticleUsecase(repo, policy)

			_, err := uc.AdminUnpublishArticle(context.Background(), tc.userID, tc.userRole, tc.articleID)
			if err != tc.expectErr {
				t.Fatalf("expected error %v, got %v", tc.expectErr, err)
			}
			repo.AssertExpectations(t)
			policy.AssertExpectations(t)
		})
	}
}
