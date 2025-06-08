import { test, expect } from '@playwright/test';

test.describe('CertMon E2E', () => {
  test('首頁顯示正確，能導向監控列表', async ({ page }) => {
    await page.goto('http://localhost:5173/');
    await expect(page.getByRole('heading', { name: /CertMon/ })).toBeVisible();
    await page.getByRole('button', { name: /監控列表/ }).click();
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
  });

  test('監控列表空狀態、新增 domain、可點擊詳細頁並返回', async ({ page }) => {
    // 1. 進入監控列表，應看到空狀態提示
    await page.goto('http://localhost:5173/domains');
    await expect(page.getByText('這裡會顯示你正在守護的網站，趕快加入第一個 domain 吧！')).toBeVisible();

    // 2. 回首頁新增一個 domain
    await page.goto('http://localhost:5173/');
    await page.getByPlaceholder('輸入你的 domain，讓我來守護！').fill('github.com');
    await page.getByRole('button', { name: /新增監控/ }).click();

    // 3. 新增後自動導向監控列表，應看到監控列表 heading
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
    // 4. 點擊第一個詳細
    await page.getByRole('button', { name: /詳細/ }).first().click();
    await expect(page.getByRole('heading', { name: /憑證詳細/ })).toBeVisible();
    // 5. 返回
    await page.getByRole('button', { name: /返回監控列表/ }).click();
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
  });
});
