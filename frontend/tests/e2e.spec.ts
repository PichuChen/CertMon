import { test, expect } from '@playwright/test';

test.describe('CertMon E2E', () => {
  test('首頁顯示正確，能導向監控列表', async ({ page }) => {
    await page.goto('http://localhost:5173/');
    await expect(page.getByRole('heading', { name: /CertMon/ })).toBeVisible();
    await page.getByRole('button', { name: /監控列表/ }).click();
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
  });

  test('監控列表顯示正確', async ({ page }) => {
    await page.goto('http://localhost:5173/domains');
    await expect(page.getByRole('heading', { name: /監控列表/ })).toBeVisible();
    await expect(page.getByRole('button', { name: /新增監控/ })).toBeVisible();
    await expect(page.getByRole('table')).toBeVisible();
  });


});
