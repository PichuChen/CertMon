import { test, expect } from '@playwright/test';

test.describe('CertMon E2E', () => {
  test('首頁顯示正確，能導向監控列表', async ({ page }) => {
    await page.goto('http://localhost:5173/');
    await expect(page.getByRole('heading', { name: /CertMon/ })).toBeVisible();
    await page.getByRole('button', { name: /監控列表/ }).click();
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
  });

  test('監控列表可點擊進入詳細頁並返回', async ({ page }) => {
    await page.goto('http://localhost:5173/domains');
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
    // 點擊第一個詳細
    await page.getByRole('button', { name: /詳細/ }).first().click();
    await expect(page.getByRole('heading', { name: /憑證詳細/ })).toBeVisible();
    // 返回
    await page.getByRole('button', { name: /返回監控列表/ }).click();
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
  });
});
